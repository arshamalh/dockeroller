package msgs

const (
	WelcomeMessage = `
Hello {name}, 
welcome to your bot,
You can use dockeroller to manage your docker daemon through different Messengers
e.g. list your images or containers:
`

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
