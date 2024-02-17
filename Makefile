DOCKER_CMD?=docker

fmt:
	go fmt ./...
.PHONY: fmt

lint:
	$(DOCKER_CMD) run --rm -v $(shell pwd):/app -w /app golangci/golangci-lint:latest run -E goimports -E godot --timeout 10m
.PHONY: lint

lint-dev:
	golangci-lint run -E goimports -E godot --timeout 10m
.PHONY: lint-dev
