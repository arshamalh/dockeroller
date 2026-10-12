package entities_test

import (
	"encoding/json"
	"testing"

	"github.com/arshamalh/dockeroller/entities"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// A trimmed sample as the Docker API sends it (cgroup v2 host).
const statsSample = `{
  "read": "2026-10-12T15:48:01.5Z",
  "preread": "2026-10-12T15:48:00.5Z",
  "pids_stats": {"current": 3},
  "cpu_stats": {
    "cpu_usage": {"total_usage": 400000000},
    "system_cpu_usage": 20000000000,
    "online_cpus": 2
  },
  "precpu_stats": {
    "cpu_usage": {"total_usage": 300000000},
    "system_cpu_usage": 18000000000,
    "online_cpus": 2
  },
  "memory_stats": {
    "usage": 52428800,
    "limit": 1048576000,
    "stats": {"inactive_file": 10485760, "file": 12582912}
  }
}`

func TestStatsFromDockerJSON(t *testing.T) {
	var s entities.Stats
	require.NoError(t, json.Unmarshal([]byte(statsSample), &s))

	// (0.1e9 / 2e9) * 2 CPUs * 100 = 10 %
	assert.InDelta(t, 10.0, s.CPUUsage(), 1e-9)
	assert.Equal(t, uint64(52428800-10485760), s.UsedMemory())
	assert.InDelta(t, 4.0, s.MemoryUsage(), 1e-9) // 40 MiB of 1000 MiB
	assert.Equal(t, 3, s.PIDsStats["current"])
}

func TestStatsCPUUsageEdgeCases(t *testing.T) {
	var first entities.Stats // Docker's first sample has an empty precpu_stats
	first.CPU.Usage.Total = 5e9
	first.CPU.SystemCPUUsage = 9e11
	first.CPU.OnlineCPUs = 4
	assert.Zero(t, first.CPUUsage(), "no previous sample: not a cumulative average")

	idle := entities.Stats{CPU: first.CPU, PreCPU: first.CPU}
	assert.Zero(t, idle.CPUUsage(), "unchanged counters: 0, not NaN")

	perCPU := entities.Stats{PreCPU: first.CPU, CPU: first.CPU}
	perCPU.CPU.Usage.Total += 1e8
	perCPU.CPU.SystemCPUUsage += 1e9
	perCPU.CPU.OnlineCPUs = 0
	perCPU.CPU.Usage.PerCPU = []uint64{1, 2}
	assert.InDelta(t, 20.0, perCPU.CPUUsage(), 1e-9, "falls back to the per-CPU list length")
}

func TestStatsMemoryCgroupV1AndNoLimit(t *testing.T) {
	v1 := entities.Stats{Memory: entities.MemoryStats{
		Usage: 100, Limit: 1000,
		Stats: map[string]uint64{"total_inactive_file": 30, "cache": 50},
	}}
	assert.Equal(t, uint64(70), v1.UsedMemory())
	assert.InDelta(t, 7.0, v1.MemoryUsage(), 1e-9)

	unknown := entities.Stats{Memory: entities.MemoryStats{Usage: 100}}
	assert.Equal(t, uint64(100), unknown.UsedMemory())
	assert.Zero(t, unknown.MemoryUsage(), "no limit reported: 0 instead of +Inf")
}
