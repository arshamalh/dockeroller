package testkit

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"strconv"
	"strings"
	"sync"
	"testing"
	"time"
)

// Call is one Bot API request the bot made.
type Call struct {
	Method string
	Params map[string]any
	// Error is the error Telegram answered with, if any.
	Error string
}

// Text returns the call's text parameter.
func (c Call) Text() string { s, _ := c.Params["text"].(string); return s }

// Message is the current state of a message the bot sent, as the user sees
// it in the chat.
type Message struct {
	ID      int
	Text    string
	Buttons [][]Button
	Deleted bool
}

// Button is an inline keyboard button.
type Button struct {
	Text string `json:"text"`
	Data string `json:"callback_data"`
}

// Telegram is a fake Bot API. Like the real one, it rejects MarkdownV2 text
// that doesn't parse and edits that change nothing.
type Telegram struct {
	t   testing.TB
	srv *httptest.Server

	mu       sync.Mutex
	calls    []Call
	messages map[int]*Message
	nextID   int
}

// NewTelegram starts a fake Bot API server for the test.
func NewTelegram(t testing.TB) *Telegram {
	tg := &Telegram{t: t, messages: map[int]*Message{}, nextID: 100}
	tg.srv = httptest.NewServer(http.HandlerFunc(tg.serve))
	t.Cleanup(tg.srv.Close)
	return tg
}

// URL is the Bot API URL to configure the bot with.
func (tg *Telegram) URL() string { return tg.srv.URL }

func (tg *Telegram) serve(w http.ResponseWriter, r *http.Request) {
	method := r.URL.Path[strings.LastIndex(r.URL.Path, "/")+1:]
	body, _ := io.ReadAll(r.Body)
	params := map[string]any{}
	_ = json.Unmarshal(body, &params)

	reply := func(result any) {
		data, _ := json.Marshal(map[string]any{"ok": true, "result": result})
		w.Write(data)
	}
	fail := func(desc string) {
		tg.mu.Lock()
		tg.calls[len(tg.calls)-1].Error = desc
		tg.mu.Unlock()
		data, _ := json.Marshal(map[string]any{"ok": false, "error_code": 400, "description": "Bad Request: " + desc})
		w.Write(data)
	}

	if method == "getUpdates" {
		time.Sleep(10 * time.Millisecond)
		reply([]any{})
		return
	}

	tg.mu.Lock()
	tg.calls = append(tg.calls, Call{Method: method, Params: params})
	tg.mu.Unlock()

	switch method {
	case "getMe":
		reply(map[string]any{"id": 1, "is_bot": true, "first_name": "Dockeroller", "username": "dockeroller_bot"})
	case "sendMessage", "editMessageText":
		text, _ := params["text"].(string)
		if params["parse_mode"] == "MarkdownV2" {
			if err := ValidateMarkdownV2(text); err != "" {
				fail(err)
				return
			}
		}
		buttons := parseButtons(params["reply_markup"])
		tg.mu.Lock()
		var id int
		if method == "sendMessage" {
			tg.nextID++
			id = tg.nextID
			tg.messages[id] = &Message{ID: id, Text: text, Buttons: buttons}
		} else {
			id = toInt(params["message_id"])
			m := tg.messages[id]
			switch {
			case m == nil || m.Deleted:
				tg.mu.Unlock()
				fail("message to edit not found")
				return
			case m.Text == text && sameButtons(m.Buttons, buttons):
				tg.mu.Unlock()
				fail("message is not modified: specified new message content and reply markup are exactly the same as a current content and reply markup of the message")
				return
			}
			m.Text, m.Buttons = text, buttons
		}
		tg.mu.Unlock()
		reply(map[string]any{"message_id": id, "chat": map[string]any{"id": toInt(params["chat_id"]), "type": "private"}, "text": text})
	case "editMessageReplyMarkup":
		id := toInt(params["message_id"])
		tg.mu.Lock()
		if m := tg.messages[id]; m != nil {
			m.Buttons = parseButtons(params["reply_markup"])
		}
		tg.mu.Unlock()
		reply(map[string]any{"message_id": id, "chat": map[string]any{"id": 1}})
	case "deleteMessage":
		tg.mu.Lock()
		if m := tg.messages[toInt(params["message_id"])]; m != nil {
			m.Deleted = true
		}
		tg.mu.Unlock()
		reply(true)
	default: // answerCallbackQuery, setMyCommands, ...
		reply(true)
	}
}

func toInt(v any) int {
	switch v := v.(type) {
	case float64:
		return int(v)
	case string:
		n, _ := strconv.Atoi(v)
		return n
	}
	return 0
}

func parseButtons(markup any) [][]Button {
	var raw []byte
	switch m := markup.(type) {
	case string:
		raw = []byte(m)
	case nil:
		return nil
	default:
		raw, _ = json.Marshal(m)
	}
	var kb struct {
		InlineKeyboard [][]Button `json:"inline_keyboard"`
	}
	_ = json.Unmarshal(raw, &kb)
	return kb.InlineKeyboard
}

func sameButtons(a, b [][]Button) bool {
	x, _ := json.Marshal(a)
	y, _ := json.Marshal(b)
	return string(x) == string(y)
}

// Calls returns every request made so far (getUpdates excluded).
func (tg *Telegram) Calls() []Call {
	tg.mu.Lock()
	defer tg.mu.Unlock()
	return append([]Call(nil), tg.calls...)
}

// Mark returns a position in the call log, for CallsSince.
func (tg *Telegram) Mark() int {
	tg.mu.Lock()
	defer tg.mu.Unlock()
	return len(tg.calls)
}

// CallsSince returns the calls made after mark.
func (tg *Telegram) CallsSince(mark int) []Call {
	tg.mu.Lock()
	defer tg.mu.Unlock()
	return append([]Call(nil), tg.calls[mark:]...)
}

// Rejected returns the calls Telegram answered with an error other than
// "not modified" (which bots legitimately hit). A test can assert it is
// empty to prove every message was valid.
func (tg *Telegram) Rejected() []Call {
	var out []Call
	for _, c := range tg.Calls() {
		if c.Error != "" && !strings.Contains(c.Error, "not modified") {
			out = append(out, c)
		}
	}
	return out
}

// Toasts returns the texts of the callback answers made after mark.
func (tg *Telegram) Toasts(mark int) []string {
	var out []string
	for _, c := range tg.CallsSince(mark) {
		if c.Method == "answerCallbackQuery" {
			s, _ := c.Params["text"].(string)
			out = append(out, s)
		}
	}
	return out
}

// Message returns the current state of a message.
func (tg *Telegram) Message(id int) Message {
	tg.mu.Lock()
	defer tg.mu.Unlock()
	if m := tg.messages[id]; m != nil {
		return *m
	}
	return Message{}
}

// Latest returns the newest message that isn't deleted.
func (tg *Telegram) Latest() Message {
	tg.mu.Lock()
	defer tg.mu.Unlock()
	best := 0
	for id, m := range tg.messages {
		if !m.Deleted && id > best {
			best = id
		}
	}
	if best == 0 {
		return Message{}
	}
	return *tg.messages[best]
}

// FindButton returns the newest message carrying a button with this label
// and that button.
func (tg *Telegram) FindButton(label string) (Message, Button, bool) {
	tg.mu.Lock()
	defer tg.mu.Unlock()
	best := 0
	var btn Button
	for id, m := range tg.messages {
		if m.Deleted || id < best {
			continue
		}
		for _, row := range m.Buttons {
			for _, b := range row {
				if b.Text == label {
					best, btn = id, b
				}
			}
		}
	}
	if best == 0 {
		return Message{}, Button{}, false
	}
	return *tg.messages[best], btn, true
}

// Labels flattens a message's buttons into their labels, row by row.
func (m Message) Labels() [][]string {
	out := make([][]string, 0, len(m.Buttons))
	for _, row := range m.Buttons {
		labels := make([]string, 0, len(row))
		for _, b := range row {
			labels = append(labels, b.Text)
		}
		out = append(out, labels)
	}
	return out
}

// String is a debugging aid.
func (m Message) String() string { return fmt.Sprintf("#%d %q %v", m.ID, m.Text, m.Labels()) }

// Eventually polls cond until it holds or the timeout passes.
func Eventually(t testing.TB, timeout time.Duration, cond func() bool, msg string, args ...any) {
	t.Helper()
	deadline := time.Now().Add(timeout)
	for time.Now().Before(deadline) {
		if cond() {
			return
		}
		time.Sleep(10 * time.Millisecond)
	}
	t.Fatalf(msg, args...)
}
