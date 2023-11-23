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

type Stats struct {
	Read      time.Time
	Preread   time.Time
	PIDsStats map[string]int `json:"pids_stats"`
	CPU       CPUStats       `json:"cpu_stats"`
	PerCPU    CPUStats       `json:"percpu_stats"`
	Memory    MemoryStats    `json:"memory_stats"`
}

// TODO: Other fields should be implemented
