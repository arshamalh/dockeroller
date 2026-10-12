package entities

import "time"

// CPUStats is the cpu_stats / precpu_stats part of a Docker stats sample.
type CPUStats struct {
	Usage struct {
		Total  uint64   `json:"total_usage"`
		PerCPU []uint64 `json:"percpu_usage"`
	} `json:"cpu_usage"`
	SystemCPUUsage uint64 `json:"system_cpu_usage"`
	OnlineCPUs     int    `json:"online_cpus"`
}

// MemoryStats is the memory_stats part of a Docker stats sample.
type MemoryStats struct {
	Usage    uint64 `json:"usage"`
	MaxUsage uint64 `json:"max_usage"`
	Limit    uint64 `json:"limit"`
	// Stats holds the raw cgroup counters. Their names differ between
	// cgroup v1 (total_inactive_file, cache, ...) and v2 (inactive_file,
	// file, ...).
	Stats map[string]uint64 `json:"stats"`
}

// Stats is one sample of GET /containers/{id}/stats. Every sample carries
// the previous one's CPU counters in PreCPU, so a single sample is enough
// to compute the current CPU usage.
type Stats struct {
	Read      time.Time      `json:"read"`
	PreRead   time.Time      `json:"preread"`
	PIDsStats map[string]int `json:"pids_stats"`
	CPU       CPUStats       `json:"cpu_stats"`
	PreCPU    CPUStats       `json:"precpu_stats"`
	Memory    MemoryStats    `json:"memory_stats"`
}

// CPUUsage returns the CPU usage in percent of one core (so 250 means two
// and a half cores busy), computed as the docker CLI does:
//
//	cpu_delta    = cpu_stats.cpu_usage.total_usage - precpu_stats.cpu_usage.total_usage
//	system_delta = cpu_stats.system_cpu_usage   - precpu_stats.system_cpu_usage
//	usage %      = cpu_delta / system_delta * online_cpus * 100
//
// It is 0 when there is no previous sample yet (the first one Docker sends)
// or the counters didn't move, instead of NaN or a cumulative average.
// The figures are relative to the host, so on Docker Desktop they match
// `docker stats` but not the VM's own task manager.
func (s Stats) CPUUsage() float64 {
	if s.PreCPU.SystemCPUUsage == 0 {
		return 0 // first sample: no previous reading to compare with
	}
	if s.CPU.Usage.Total <= s.PreCPU.Usage.Total || s.CPU.SystemCPUUsage <= s.PreCPU.SystemCPUUsage {
		return 0
	}
	cpuDelta := float64(s.CPU.Usage.Total - s.PreCPU.Usage.Total)
	systemDelta := float64(s.CPU.SystemCPUUsage - s.PreCPU.SystemCPUUsage)
	cpus := s.CPU.OnlineCPUs
	if cpus == 0 {
		cpus = len(s.CPU.Usage.PerCPU)
	}
	return cpuDelta / systemDelta * float64(cpus) * 100
}

// UsedMemory returns the memory in use without the page cache, the number
// `docker stats` shows: usage minus inactive file pages (total_inactive_file
// on cgroup v1, inactive_file on v2).
func (s Stats) UsedMemory() uint64 {
	m := s.Memory
	if v, ok := m.Stats["total_inactive_file"]; ok && v < m.Usage {
		return m.Usage - v
	}
	if v, ok := m.Stats["inactive_file"]; ok && v < m.Usage {
		return m.Usage - v
	}
	return m.Usage
}

// MemoryUsage returns UsedMemory as a percentage of the limit, or 0 when
// the daemon reports no limit.
func (s Stats) MemoryUsage() float64 {
	if s.Memory.Limit == 0 {
		return 0
	}
	return float64(s.UsedMemory()) / float64(s.Memory.Limit) * 100
}
