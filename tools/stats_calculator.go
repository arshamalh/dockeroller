package tools

import (
	"github.com/arshamalh/dockeroller/models"
)

// These are formulas to calculate stats
// cpu_delta = cpu_stats.cpu_usage.total_usage - precpu_stats.cpu_usage.total_usage
// system_cpu_delta = cpu_stats.system_cpu_usage - precpu_stats.system_cpu_usage
// number_of_cpu_cores = lenght(cpu_stats.cpu_usage.percpu_usage) or cpu_stats.online_cpus
// CPU usage % = (cpu_delta / system_cpu_delta) * number_of_cpu_cores * 100.0
//
// used_memory = memory_stats.usage - memory_stats.stats.cache
// available_memory = memory_stats.limit
// Memory usage % = (used_memory / available_memory) * 100.0
//
// CPU usage are based on the host, not VM, so in Windows devices, you may see different cpu usages in the Bot and Docker desktop and docker cli,
// but, that's calculable. both are correct anyway.
func StatsCalculator(stat models.Stats) (cpu_usage float64, memory_usage float64) {
	// CPU usage calculation
	cpu_delta := stat.CPU.Usage.Total - stat.PerCPU.Usage.Total
	system_cpu_delta := stat.CPU.SystemCPUUsage - stat.PerCPU.SystemCPUUsage
	number_of_cpu_cores := stat.CPU.OnlineCPUs
	if number_of_cpu_cores == 0 {
		number_of_cpu_cores = len(stat.CPU.Usage.PerCPU)
	}
	cpu_usage = (float64(cpu_delta) / float64(system_cpu_delta)) * float64(number_of_cpu_cores) * 100

	// Memory usage calculation
	used_memory := stat.Memory.Usage - stat.Memory.Stats.Cache
	memory_usage = (float64(used_memory) / float64(stat.Memory.Limit)) * 100
	return
}
