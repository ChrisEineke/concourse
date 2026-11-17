.PHONY := build

build:
	CGO_ENABLED=0 go build .

container:
	docker-buildx build -t concourse/deployer .
