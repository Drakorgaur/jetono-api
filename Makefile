#REGISTRY ?= docker.ops.iszn.cz/pomo/nats-cluster/nats
REGISTRY ?= pomo.ops.dszn.cz:30500/pomo/jetono

converge:
	werf converge --repo $(REGISTRY) --insecure-registry --dev

compose-up:
	werf compose up \
		--repo $(REGISTRY) --dev

render:
	werf render \
		--repo $(REGISTRY) \
		--dev \
		--set environment=dev,locality=ko \
		--release nats-cluster-dev-ko

bundle-render:
	werf bundle render \
		--repo $(REGISTRY)
