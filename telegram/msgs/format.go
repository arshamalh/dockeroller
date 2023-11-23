package msgs

import (
	"fmt"
	"strings"

	"github.com/arshamalh/dockeroller/entities"
	"github.com/arshamalh/dockeroller/tools"
)

// Monospace font is enabled by ` charater that is not supported in Go multiline string literals
// So we should format it using ” and replacing them.
// characters ()_-.>=< are also reserved by telegram and we will replace them with their escaped ones.
func FmtMono(input string) string {
	return strings.NewReplacer(
		"''", "`",
		"(", "\\(",
		")", "\\)",
		"_", "\\_",
		".", "\\.",
		"-", "\\-",
		"=", "\\=",
		"<", "\\<",
		">", "\\>",
	).Replace(input)
}

func FmtContainer(container *entities.Container) string {
	response := strings.NewReplacer(
		"{name}", container.Name,
		"{image}", container.Image,
		"{status}", container.Status,
	).Replace(Container)
	response = FmtMono(response)
	return response
}

func FmtImage(image *entities.Image) string {
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

func FmtStats(stat entities.Stats) string {
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

func FmtContainerRenamed(old_name, new_name string) string {
	response := strings.NewReplacer(
		"{old_name}", old_name,
		"{new_name}", new_name,
	).Replace(ContainerRenamed)
	response = FmtMono(response)
	return response
}

func FmtImageTagged(id, tag string) string {
	response := strings.NewReplacer(
		"{id}", id,
		"{tag}", tag,
	).Replace(ImageTagged)
	response = FmtMono(response)
	return response
}
