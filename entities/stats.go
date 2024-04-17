package entities

import "time"

type CPUStats struct {
	Usage struct {
		Total  uint   `json:"total_usage"`
		PerCPU []uint `json:"percpu_usage"`
	} `json:"cpu_usage"`
	SystemCPUUsage uint `json:"system_cpu_usage"`
	OnlineCPUs     int  `json:"online_cpus"`
}

type MemoryStats struct {
	Usage    int64 `json:"usage"`
	MaxUsage int64 `json:"max_usage"`
	Limit    int64
	Stats    struct {
		Cache int64
	} `json:"stats"`
}

// TODO: Other fields should be implemented

type Stats struct {
	Read      time.Time
	Preread   time.Time
	PIDsStats map[string]int `json:"pids_stats"`
	CPU       CPUStats       `json:"cpu_stats"`
	PerCPU    CPUStats       `json:"percpu_stats"`
	Memory    MemoryStats    `json:"memory_stats"`
}

// These are formulas to calculate stats
// cpu_delta = cpu_stats.cpu_usage.total_usage - precpu_stats.cpu_usage.total_usage
// system_cpu_delta = cpu_stats.system_cpu_usage - precpu_stats.system_cpu_usage
// number_of_cpu_cores = lenght(cpu_stats.cpu_usage.percpu_usage) or cpu_stats.online_cpus
// CPU usage % = (cpu_delta / system_cpu_delta) * number_of_cpu_cores * 100.0
//
// CPU usage are based on the host, not VM, so in Windows devices, you may see different cpu usages in the Bot and Docker desktop and docker cli,
// but, that's calculable. both are correct anyway.
func (stat Stats) CPUUsage() float64 {
	cpu_delta := stat.CPU.Usage.Total - stat.PerCPU.Usage.Total
	system_cpu_delta := stat.CPU.SystemCPUUsage - stat.PerCPU.SystemCPUUsage
	number_of_cpu_cores := stat.CPU.OnlineCPUs
	if number_of_cpu_cores == 0 {
		number_of_cpu_cores = len(stat.CPU.Usage.PerCPU)
	}
	cpu_usage := (float64(cpu_delta) / float64(system_cpu_delta)) * float64(number_of_cpu_cores) * 100
	return cpu_usage
}

// Simplifying the memory usage calculations
// used_memory = memory_stats.usage - memory_stats.stats.cache
// available_memory = memory_stats.limit
// Memory usage % = (used_memory / available_memory) * 100.0
func (stat Stats) MemoryUsage() float64 {
	used_memory := stat.Memory.Usage - stat.Memory.Stats.Cache
	memory_usage := (float64(used_memory) / float64(stat.Memory.Limit)) * 100
	return memory_usage
}
