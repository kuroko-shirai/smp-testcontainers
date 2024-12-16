#/bin/bash
# Use Bash instead of SH
export SHELL := /bin/bash

.DEFAULT_GOAL := controll

GOPATH := $(shell go env GOPATH)

SERVER_PATH := cmd/app

# Run the docker
.PHONY: up
up:
	@echo "Up"
	@docker-compose up -d

# Run the docker
.PHONY: down
down:
	@echo "Down"
	@docker-compose down

# Run the server
.PHONY: server
server:
	@echo "Server running..."
	@go run -race $(SERVER_PATH)/main.go
