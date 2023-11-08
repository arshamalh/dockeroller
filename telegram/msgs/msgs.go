package msgs

const (
	Container = `
''Container name: '' {name},
''Used Image:     '' {image},
''Status:         '' {status}
`
	Image = `
''ID:   '' {id},
''Size: '' {size},
''Tags: '' {tags}
`

	Stat = `
''CPU Usage:    '' {cpu_usage} %
''Online CPUs:  '' {online_cpus}
''Memory Usage: '' {memory_usage} ({memory_usage%} %)
''Avaiable Mem: '' {avaiable_memory}
`
)
