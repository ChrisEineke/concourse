.PHONY := build container test
.DEFAULT_GOAL := build

build:
	$(MAKE) -C cmd/concourse-swarm-key-maker build
	$(MAKE) -C cmd/deployer build
	$(MAKE) -C cmd/node-orchestrator build

container:
	$(MAKE) -C cmd/deployer container
	$(MAKE) -C cmd/node-orchestrator container

test:
	$(MAKE) -C cmd/concourse-swarm-key-maker test
	$(MAKE) -C cmd/deployer test
	$(MAKE) -C cmd/node-orchestrator test
