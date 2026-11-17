.PHONY := build test

build:
	CGO_ENABLED=0 go build .

container:
	docker-buildx build -t concourse/swarm-key-maker .

test:
	CGO_ENABLED=0 go test -v ./...
