// Package testkit has fakes for testing dockeroller end to end without
// Telegram or a Docker daemon: a Bot API server that checks messages the way
// Telegram does, and an in-memory Docker.
package testkit

import (
	"fmt"
	"strings"
)

// ValidateMarkdownV2 reports whether text would be accepted by Telegram with
// parse_mode=MarkdownV2. It returns "" when valid, otherwise an error
// description shaped like Telegram's ("can't parse entities: ...").
//
// It covers what dockeroller relies on: reserved characters outside
// entities, *bold*, _italic_, __underline__, ~strike~, ||spoiler||, `code`,
// ```pre``` and [links](...). Inside code and pre only '`' and '\' are
// special, as in Telegram.
func ValidateMarkdownV2(text string) string {
	const reserved = "_*[]()~`>#+-=|{}.!"
	type open struct {
		marker string
		offset int
	}
	var stack []open
	top := func() string {
		if len(stack) == 0 {
			return ""
		}
		return stack[len(stack)-1].marker
	}
	for i := 0; i < len(text); {
		c := text[i]
		if c == '\\' {
			if i+1 >= len(text) {
				return "can't parse entities: Character '\\' is reserved and must be escaped with the preceding '\\'"
			}
			i += 2
			continue
		}
		if t := top(); t == "`" || t == "```" {
			if strings.HasPrefix(text[i:], t) {
				stack = stack[:len(stack)-1]
				i += len(t)
				continue
			}
			i++
			continue
		}
		switch {
		case strings.HasPrefix(text[i:], "```"):
			stack = append(stack, open{"```", i})
			i += 3
			continue
		case c == '`':
			stack = append(stack, open{"`", i})
			i++
			continue
		case strings.HasPrefix(text[i:], "||"), strings.HasPrefix(text[i:], "__"):
			m := text[i : i+2]
			if top() == m {
				stack = stack[:len(stack)-1]
			} else {
				stack = append(stack, open{m, i})
			}
			i += 2
			continue
		case c == '*' || c == '_' || c == '~':
			m := string(c)
			if top() == m {
				stack = stack[:len(stack)-1]
			} else {
				for _, o := range stack {
					if o.marker == m {
						return fmt.Sprintf("can't parse entities: Can't find end of the entity starting at byte offset %d", o.offset)
					}
				}
				stack = append(stack, open{m, i})
			}
			i++
			continue
		case c == '[':
			j := strings.Index(text[i:], "](")
			if j < 0 {
				return "can't parse entities: Character '[' is reserved and must be escaped with the preceding '\\'"
			}
			k := strings.IndexByte(text[i+j+2:], ')')
			if k < 0 {
				return "can't parse entities: Character '(' is reserved and must be escaped with the preceding '\\'"
			}
			i += j + 2 + k + 1
			continue
		case strings.IndexByte(reserved, c) >= 0:
			return fmt.Sprintf("can't parse entities: Character '%c' is reserved and must be escaped with the preceding '\\'", c)
		}
		i++
	}
	if len(stack) > 0 {
		return fmt.Sprintf("can't parse entities: Can't find end of the entity starting at byte offset %d", stack[len(stack)-1].offset)
	}
	return ""
}
