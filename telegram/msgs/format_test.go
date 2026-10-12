package msgs_test

import (
	"strings"
	"testing"
	"unicode/utf8"

	"github.com/arshamalh/dockeroller/entities"
	"github.com/arshamalh/dockeroller/internal/testkit"
	"github.com/arshamalh/dockeroller/telegram/msgs"
	"github.com/stretchr/testify/assert"
)

func TestShortImageID(t *testing.T) {
	assert.Equal(t, "863dd03b5190", msgs.ShortImageID("sha256:863dd03b51905b838bbe0e32d14295bbe855412b2b23c74b4061551c1d071236"))
	assert.Equal(t, "abc", msgs.ShortImageID("abc"))
}

func TestImageArgs(t *testing.T) {
	args := msgs.ImageArgs(&entities.Image{ID: "sha256:0123456789abcdef", Size: 2_500_000, Status: entities.ImageStatusUnUsedDangling, CreatedAt: "2026-01-02 03:04:05"})
	assert.Equal(t, []string{
		"id", "0123456789ab",
		"size", "2.50 MegaBytes",
		"tags", "<none>",
		"status", "Dangling",
		"created", "2026-01-02 03:04:05",
	}, args)

	args = msgs.ImageArgs(&entities.Image{Tags: []string{"a:1", "b:2"}})
	assert.Equal(t, "a:1, b:2", args[5])
}

func TestLogs(t *testing.T) {
	text := msgs.Logs("worker_2.prod", []string{"a_b (x) [y] *z*", "back`tick \\ slash"})
	assert.Equal(t, "🪵 *worker\\_2\\.prod*\n```\na_b (x) [y] *z*\nback\\`tick \\\\ slash\n```", text)
	assert.Empty(t, testkit.ValidateMarkdownV2(text))

	assert.Contains(t, msgs.Logs("x", nil), "(no output yet)")
	assert.Contains(t, msgs.Logs("x", []string{"  ", ""}), "(no output yet)")
}

func TestLogsTruncation(t *testing.T) {
	// Multi-byte text (Farsi, emoji) in a long log: the cut must not split a
	// character, or Telegram rejects the message as invalid UTF-8. The
	// padding shifts where the cut lands through every byte alignment.
	line := "سلام دنیا 🐳 " + strings.Repeat("ب", 50)
	for pad := 0; pad < 4; pad++ {
		var many []string
		for len(strings.Join(many, "\n")) < 10_000 {
			many = append(many, line)
		}
		last := line + strings.Repeat("x", pad) // moves the cut, which is measured from the end
		many = append(many, last)
		text := msgs.Logs("svc", many)

		assert.True(t, utf8.ValidString(text), "pad %d: truncated logs must stay valid UTF-8", pad)
		assert.Contains(t, text, "```\n…", "the cut is marked")
		assert.Less(t, len(text), 4096, "fits in a Telegram message")
		assert.True(t, strings.HasSuffix(text, last+"\n```"), "the newest line is kept whole")
		assert.Empty(t, testkit.ValidateMarkdownV2(text))
	}
}

func TestStats(t *testing.T) {
	var s entities.Stats
	s.CPU.Usage.Total, s.CPU.SystemCPUUsage, s.CPU.OnlineCPUs = 400, 2000, 2
	s.PreCPU.Usage.Total, s.PreCPU.SystemCPUUsage = 300, 1000
	s.Memory.Usage, s.Memory.Limit = 3_000_000, 1_000_000_000

	text := msgs.Stats("db.primary", s)
	assert.Equal(t, "📊 *db\\.primary*\n"+
		"`CPU usage:   ` 20\\.00 %\n"+
		"`Online CPUs: ` 2\n"+
		"`Memory:      ` 3\\.00 MegaBytes \\(0\\.30 %\\)\n"+
		"`Mem limit:   ` 1\\.00 GigaBytes", text)
	assert.Empty(t, testkit.ValidateMarkdownV2(text))
}

func TestEveryTranslationRenders(t *testing.T) {
	for key, tpl := range msgs.EN {
		if strings.HasPrefix(key, "toast.") || strings.HasPrefix(key, "btn.") || key == "lang" {
			continue // plain text: toasts and button labels are not parsed
		}
		// Rendering escapes the template's punctuation but keeps its
		// markup; simulate that and check Telegram would accept it.
		rendered := strings.NewReplacer("{", "", "}", "").Replace(tpl)
		rendered = escapeBasics(rendered)
		assert.Empty(t, testkit.ValidateMarkdownV2(rendered), "template %q", key)
	}
}

// escapeBasics mirrors tbb.EscapeMarkdownV2Basics without importing tbb's
// internals into the assertion.
func escapeBasics(s string) string {
	r := strings.NewReplacer("[", `\[`, "]", `\]`, "(", `\(`, ")", `\)`, "~", `\~`, ">", `\>`, "<", `\<`,
		"#", `\#`, "+", `\+`, "-", `\-`, "=", `\=`, "|", `\|`, "{", `\{`, "}", `\}`, ".", `\.`, "!", `\!`)
	return r.Replace(s)
}
