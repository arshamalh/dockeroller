package msgs

import "strings"

const (
	Container = `
''Container name: '' {Name},
''Used Image:     '' {Image},
''Status:         '' {Status},
`
)

// Monospace font is enabled by ` charater that is not supported in Go multiline string literals
// So we should format it using '' and replacing them.
func FmtMono(input string) string {
	return strings.NewReplacer(
		"''", "`",
	).Replace(input)
}
