.PHONY := build container test

build:
	CGO_ENABLED=0 go build .

container:
	docker-buildx build -t concourse/deployer .

test:
	CGO_ENABLED=0 go test -v ./...
