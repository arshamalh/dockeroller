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
		"{size}", tools.SizeToHumanReadable(image.Size),
		"{tags}", fmt.Sprint(image.Tags),
		"{status}", fmt.Sprint(image.Status),
		"{created_at}", fmt.Sprint(image.CreatedAt),
	).Replace(Image)
	response = FmtMono(response)
	return response
}

func FmtStats(stat models.Stats) string {
	cpu_usage, memory_usage_percent := tools.StatsCalculator(stat)
	response := strings.NewReplacer(
		"{cpu_usage}", fmt.Sprintf("%.2f", cpu_usage),
		"{online_cpus}", fmt.Sprint(stat.CPU.OnlineCPUs),
		"{memory_usage}", tools.SizeToHumanReadable(stat.Memory.Usage),
		"{memory_usage%}", fmt.Sprintf("%.2f", memory_usage_percent),
		"{avaiable_memory}", tools.SizeToHumanReadable(stat.Memory.Limit),
	).Replace(Stat)
	response = FmtMono(response)
	return response
}
