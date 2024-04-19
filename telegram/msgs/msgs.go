package msgs

const (
	WelcomeMessage = `
Hello {name}, 
welcome to your bot,
You can use dockeroller to manage your docker daemon through different Messengers
e.g. list your images or containers:
`

	ContainerRenamed = `
Container {old_name} successfully renamed to {new_name}
`

	ImageTagged = `
Image {id} successfully tagged {tag}
`

	ContainerNewNameInput = `
Enter new container name:	
`

	ImageNewNameInput = `
Enter new image tag:
`

	InvalidImageTag = `
invalid tag, tags should be:
- all lower case
- with no space
`

	Container = `
''Name:   '' {name},
''Image:  '' {image},
''Status: '' {status}
`

	Image = `
''ID:      '' ''{id1}''...{id2},
''Size:    '' {size},
''Tags:    '' {tags}
''Status:  '' {status}
''Created: '' {created_at}
`

	Stat = `
''CPU Usage:    '' {cpu_usage} %
''Online CPUs:  '' {online_cpus}
''Memory Usage: '' {memory_usage} ({memory_usage%} %)
''Available Mem: '' {available_memory}
`
)
