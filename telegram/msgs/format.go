package msgs

import (
	"fmt"
	"strings"
	"unicode/utf8"

	"github.com/arshamalh/dockeroller/entities"
	"github.com/arshamalh/dockeroller/tools"
	"github.com/arshamalh/tbb"
)

// ContainerArgs returns the template values of the Container message.
func ContainerArgs(c *entities.Container) []string {
	return []string{
		"name", c.Name,
		"image", c.Image,
		"state", string(c.State),
		"status", c.Status,
	}
}

// ShortImageID turns "sha256:863dd03b5190…1236" into "863dd03b5190", the
// form `docker images` shows.
func ShortImageID(id string) string {
	id = strings.TrimPrefix(id, "sha256:")
	if len(id) > 12 {
		return id[:12]
	}
	return id
}

// ImageArgs returns the template values of the Image message.
func ImageArgs(img *entities.Image) []string {
	tags := "<none>"
	if len(img.Tags) > 0 {
		tags = strings.Join(img.Tags, ", ")
	}
	return []string{
		"id", ShortImageID(img.ID),
		"size", tools.SizeToHumanReadable(img.Size),
		"tags", tags,
		"status", string(img.Status),
		"created", img.CreatedAt,
	}
}

// Stats renders a stats sample as MarkdownV2. It is built here rather than
// through a template because it is rendered by a background stream, outside
// any handler context.
func Stats(name string, s entities.Stats) string {
	rows := [][2]string{
		{"CPU usage:   ", fmt.Sprintf("%.2f %%", s.CPUUsage())},
		{"Online CPUs: ", fmt.Sprint(s.CPU.OnlineCPUs)},
		{"Memory:      ", fmt.Sprintf("%s (%.2f %%)", tools.SizeToHumanReadable(int64(s.UsedMemory())), s.MemoryUsage())},
		{"Mem limit:   ", tools.SizeToHumanReadable(int64(s.Memory.Limit))},
	}
	var b strings.Builder
	b.WriteString("📊 *" + tbb.EscapeMarkdownV2(name) + "*\n")
	for _, r := range rows {
		b.WriteString("`" + r[0] + "` " + tbb.EscapeMarkdownV2(r[1]) + "\n")
	}
	return strings.TrimSuffix(b.String(), "\n")
}

// maxLogChars keeps log messages under Telegram's 4096-character limit,
// leaving room for the header and the code fence.
const maxLogChars = 3500

// tail returns the end of s, at most max bytes long. It starts at a line
// boundary when one is close, and never in the middle of a UTF-8 character
// (Telegram rejects invalid UTF-8).
func tail(s string, max int) string {
	start := len(s) - max
	for start < len(s) && !utf8.RuneStart(s[start]) {
		start++
	}
	if nl := strings.IndexByte(s[start:], '\n'); nl >= 0 && nl < 200 {
		start += nl + 1
	}
	return s[start:]
}

// Logs renders the last log lines as a MarkdownV2 pre block. Only '`' and '\'
// need escaping inside it; anything else would show up as literal
// backslashes.
func Logs(name string, lines []string) string {
	body := strings.Join(lines, "\n")
	if len(body) > maxLogChars {
		body = "…" + tail(body, maxLogChars)
	}
	if strings.TrimSpace(body) == "" {
		body = "(no output yet)"
	}
	return "🪵 *" + tbb.EscapeMarkdownV2(name) + "*\n```\n" + tbb.EscapeMarkdownV2Code(body) + "\n```"
}
