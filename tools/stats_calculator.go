package tools

import "github.com/arshamalh/dockeroller/models"

// These are formulas to calculate stats
// cpu_delta = cpu_stats.cpu_usage.total_usage - precpu_stats.cpu_usage.total_usage
// system_cpu_delta = cpu_stats.system_cpu_usage - precpu_stats.system_cpu_usage
// number_cpus = lenght(cpu_stats.cpu_usage.percpu_usage) or cpu_stats.online_cpus
// CPU usage % = (cpu_delta / system_cpu_delta) * number_cpus * 100.0
//
// used_memory = memory_stats.usage - memory_stats.stats.cache
// available_memory = memory_stats.limit
// Memory usage % = (used_memory / available_memory) * 100.0
//
func StatsCalculator(stat models.Stats) (cpu_usage uint, memory_usage uint) {
	// CPU usage calculation
	cpu_delta := stat.CPU.Usage.Total - stat.PerCPU.Usage.Total
	system_cpu_delta := stat.CPU.SystemCPUUsage - stat.PerCPU.SystemCPUUsage
	number_of_cpus := stat.CPU.OnlineCPUs
	if number_of_cpus == 0 {
		number_of_cpus = len(stat.CPU.Usage.PerCPU)
	}
	cpu_usage = cpu_delta / system_cpu_delta * uint(number_of_cpus) * 100

	// Memory usage calculation
	used_memory := stat.Memory.Usage - stat.Memory.Stats.Cache
	memory_usage = used_memory / stat.Memory.Limit * 100
	return
}
