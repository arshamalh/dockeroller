# Default values, but should be overriden like `make run File=data/another.json
version=latest
name=dockeroller

help: # Generate list of targets with descriptions
	@grep '^.*\:\s#.*' Makefile | sed 's/\(.*\) # \(.*\)/\1 \2/' | column -t -s ":"

build: # Build the binary for current local system
	go build -o dockeroller .

build-docker: # Build the docker image
	docker build -t ${name}:${version} .

sample-docker: # Run docker with sample start command
	docker run --rm ${name}:${version} start --token="something"

gen:
	go generate ./...
	