.PHONY: build
build:
	go build -v ./cmd/proxy

.DEFAULT_GOAL := build