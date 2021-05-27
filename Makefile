# Go parameters
GOCMD=go
NPMCMD=npm
NPMBUILD=$(NPMCMD) run build
GOBUILD=$(GOCMD) build
GOFMT=gofmt
BINARY_NAME=thunderdome-planning-poker
BINARY_UNIX=$(BINARY_NAME)_unix
BINARY_WINDOWS=thunderdome-planning-poker.exe
GORELEASER=goreleaser release --rm-dist
NEXT_DOCKER_TAG=stevenweathers/thunderdome-planning-poker:next
LATEST_DOCKER_TAG=stevenweathers/thunderdome-planning-poker:latest
VERSION_TAG := $(shell git tag --sort=-version:refname | head -n 1)
GOBUILDTAG=-ldflags "-X main.version=$(VERSION_TAG)"
DOCKER_BUILD_VERSION=--build-arg BUILD_VERSION=${VERSION_TAG}

all: build
build-deps: 
	$(NPMBUILD)

build: 
	$(NPMBUILD)
	$(GOBUILD) -o $(BINARY_NAME) -v

clean: 
	$(GOCLEAN)
	rm -f $(BINARY_NAME)
	rm -f $(BINARY_UNIX)
	rm -f $(BINARY_WINDOWS)
	rm -f *-packr.go
	rm -rf dist
	rm -rf release
	rm -rf pkged*.go

format:
	$(GOFMT) -s -w datasrc.go
	$(GOFMT) -s -w handlers.go
	$(GOFMT) -s -w client.go
	$(GOFMT) -s -w hub.go
	$(GOFMT) -s -w main.go
	$(GOFMT) -s -w utils.go

# Cross compilation
build-linux:
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 $(GOBUILD) -o $(BINARY_UNIX) -v

build-windows:
	CGO_ENABLED=0 GOOS=windows GOARCH=amd64 $(GOBUILD) -o $(BINARY_WINDOWS) -v

dev: 
	$(NPMBUILD)
	$(GOBUILD) -o $(BINARY_NAME) -v

	SMTP_SECURE="false" DB_HOST="localhost" APP_DOMAIN=".127.0.0.1" COOKIE_SECURE="false" ./$(BINARY_NAME) live
dev-go: 
	$(GOBUILD) -o $(BINARY_NAME) -v

	SMTP_SECURE="false" DB_HOST="localhost" APP_DOMAIN=".127.0.0.1" COOKIE_SECURE="false" ./$(BINARY_NAME) live
run:
	SMTP_SECURE="false" DB_HOST="localhost" APP_DOMAIN=".127.0.0.1" COOKIE_SECURE="false" ./$(BINARY_NAME)

gorelease:
	$(GORELEASER)

gorelease-dev-dry:
	$(GORELEASER) --skip-publish --skip-validate

gorelease-dry:
	$(GORELEASER) --skip-publish

gorelease-snapshot:
	$(GORELEASER) --snapshot

build-next-image:
	docker build ./ -f ./build/Dockerfile -t $(NEXT_DOCKER_TAG) ${DOCKER_BUILD_VERSION}

push-next-image:
	docker push $(NEXT_DOCKER_TAG)

build-latest-image:
	docker build ./ -f ./build/Dockerfile -t $(LATEST_DOCKER_TAG) ${DOCKER_BUILD_VERSION}

push-latest-image:
	docker push $(LATEST_DOCKER_TAG)