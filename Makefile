.DEFAULT_GOAL := run
SHELL := /bin/bash
APP ?= $(shell basename $$(pwd))
COMMIT_SHA = $(shell git rev-parse --short HEAD)

.PHONY: help
## help: prints this help message
help:
	@echo "Usage:"
	@sed -n 's/^##//p' ${MAKEFILE_LIST} | column -t -s ':' |  sed -e 's/^/ /'

.PHONY: run
## run: runs via air hot-reloader
run:
	source .env_private; air

.PHONY: build
## build: builds the application
build: clean
	@echo "Building binary ..."
	go build -o ${APP}

.PHONY: clean
## clean: cleans up binary files
clean:
	@echo "Cleaning up ..."
	@go clean

.PHONY: install
## install: installs the application
install:
	@echo "Installing binary ..."
	go install

.PHONY: test
## test: runs go test with the race detector
test: build
	GOARCH=amd64 GOOS=linux go test -v -race ./...

.PHONY: install-air
## install-air: installs air hot-reloader
install-air:
	go install github.com/cosmtrek/air@v1.48.0
	#go install github.com/cosmtrek/air@latest

.PHONY: init
## init: sets up go modules
init:
	@echo "Setting up modules ..."
	@go mod init 2>/dev/null; go mod tidy && go mod vendor

.PHONY: cleanup
cleanup: docker-cleanup
.PHONY: docker-cleanup
## docker-cleanup: cleans up local docker images and volumes
docker-cleanup:
	docker system prune --volumes -a

########################################################################################################################
####### docker/kubernetes related stuff ################################################################################
########################################################################################################################
.PHONY: image-login
## image-login: login to docker hub
image-login:
	@export PATH="$$HOME/bin:$$PATH"
	@echo $$DOCKER_PASS | docker login -u $$DOCKER_USER --password-stdin

.PHONY: image-build
## image-build: build docker image
image-build: build
	@export PATH="$$HOME/bin:$$PATH"
	docker build -t jamesclonk/${APP}:${COMMIT_SHA} .

.PHONY: image-publish
## image-publish: build and publish docker image
image-publish:
	@export PATH="$$HOME/bin:$$PATH"
	docker push jamesclonk/${APP}:${COMMIT_SHA}
	docker tag jamesclonk/${APP}:${COMMIT_SHA} jamesclonk/${APP}:latest
	docker push jamesclonk/${APP}:latest

.PHONY: image-run
## image-run: run docker image
image-run:
	@export PATH="$$HOME/bin:$$PATH"
	docker run --rm -p 8080:8080 jamesclonk/${APP}:${COMMIT_SHA}
