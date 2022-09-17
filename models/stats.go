package models

import "time"

type CPUUsageStats struct {
	Usage struct {
		Total  uint   `json:"total_usage"`
		PerCPU []uint `json:"percpu_usage"`
	} `json:"cpu_usage"`
	SystemCPUUsage uint `json:"system_cpu_usage"`
	OnlineCPUs     int  `json:"online_cpus"`
}

type Stats struct {
	Read      time.Time
	Preread   time.Time
	PIDsStats map[string]int `json:"pids_stats"`
	CPU       CPUUsageStats  `json:"cpu_stats"`
	PerCPU    CPUUsageStats  `json:"percpu_stats"`
	Memory    struct {
		Usage    uint
		MaxUsage uint `json:"max_usage"`
		Limit    uint
		Stats    struct {
			Cache uint
		} `json:"stats"`
	} `json:"memory_stats"`
}

// TODO: Other fields should be implemented
