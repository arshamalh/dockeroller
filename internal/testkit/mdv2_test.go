package testkit

import "testing"

func TestValidateMarkdownV2(t *testing.T) {
	valid := []string{
		"plain text",
		`*bold* _italic_ __under__ ~strike~ ||spoiler||`,
		"`code with . and - and (x)` outside\\.",
		"```\nlogs <tick> a_b (x) #1\n```",
		`escaped \. \- \( \) \! \_ \*`,
		"[a link](https://example.com/a_b)",
		"`Name:   ` web\\-1\n`Image:  ` dclr\\-fixture:latest",
	}
	for _, s := range valid {
		if err := ValidateMarkdownV2(s); err != "" {
			t.Errorf("%q: unexpected %s", s, err)
		}
	}

	invalid := map[string]string{
		"Hello john_doe":       "can't parse entities: Can't find end of the entity starting at byte offset 10",
		"version 1.2":          "can't parse entities: Character '.' is reserved and must be escaped with the preceding '\\'",
		"tags [none]":          "can't parse entities: Character '[' is reserved and must be escaped with the preceding '\\'",
		"*bold":                "can't parse entities: Can't find end of the entity starting at byte offset 0",
		"`unterminated code":   "can't parse entities: Can't find end of the entity starting at byte offset 0",
		"trailing backslash\\": "can't parse entities: Character '\\' is reserved and must be escaped with the preceding '\\'",
	}
	for s, want := range invalid {
		if got := ValidateMarkdownV2(s); got != want {
			t.Errorf("%q:\n got %q\nwant %q", s, got, want)
		}
	}
}
