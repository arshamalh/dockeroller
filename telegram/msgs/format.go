package msgs

import (
	"fmt"
	"strings"

	"github.com/arshamalh/dockeroller/models"
	"github.com/arshamalh/dockeroller/tools"
)

// Monospace font is enabled by ` charater that is not supported in Go multiline string literals
// So we should format it using ” and replacing them.
func FmtMono(input string) string {
	return strings.NewReplacer(
		"''", "`",
		"(", "\\(", // ()_-. are reserved by telegram.
		")", "\\)",
		"_", "\\_",
		".", "\\.",
		"-", "\\-",
	).Replace(input)
}

func FmtContainer(container *models.Container) string {
	response := strings.NewReplacer(
		"{name}", container.Name,
		"{image}", container.Image,
		"{status}", container.Status,
	).Replace(Container)
	response = FmtMono(response)
	return response
}

func FmtImage(image *models.Image) string {
	response := strings.NewReplacer(
		"{id}", image.ID,
		"{size}", fmt.Sprint(image.Size),
		"{tags}", fmt.Sprint(image.Tags),
	).Replace(Image)
	response = FmtMono(response)
	return response
}

func FmtStats(stat models.Stats) string {
	cpu_usage, memory_usage := tools.StatsCalculator(stat)
	response := strings.NewReplacer(
		"{cpu_usage}", fmt.Sprint(cpu_usage),
		"{online_cpus}", fmt.Sprint(stat.CPU.OnlineCPUs),
		"{memory_usage}", fmt.Sprint(memory_usage),
		"{avaiable_memory}", fmt.Sprint(stat.Memory.Limit),
	).Replace(Stat)
	response = FmtMono(response)
	return response
}
