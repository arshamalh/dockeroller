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
''CPU Usage:    '' {cpu_usage},
''Memory Usage: '' {memory_usage}
''Online CPUs:  '' {online_cpus}
`
)
