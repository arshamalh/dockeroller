# Default values, but should be overriden like `make run File=data/another.json
version=latest
name=dockeroller

help: # Generate list of targets with descriptions
	@grep '^.*\:\s#.*' Makefile | sed 's/\(.*\) # \(.*\)/\1 \2/' | column -t -s ":"

build: # Build the binary for current local system
	go build -o dockeroller .

build-docker: # Build the docker image
	docker build -t ${name}:${version} .

sample-docker: # Run the image against the host's Docker daemon (TOKEN and WHITELISTED_IDS from the environment)
	docker run --rm -v /var/run/docker.sock:/var/run/docker.sock -v dockeroller-data:/data \
		-e TOKEN -e WHITELISTED_IDS ${name}:${version} start

test: # Run the unit and end-to-end tests (no Docker daemon or Telegram needed)
	go test -race ./...

test-docker: # Also run the integration tests against the local Docker daemon
	DOCKEROLLER_DOCKER_TESTS=1 go test -race ./...

gen:
	go generate ./...
	