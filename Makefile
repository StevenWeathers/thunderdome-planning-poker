# Go parameters
GOCMD=go
STATICPACKCMD=packr
NPMCMD=npm
NPMBUILD=$(NPMCMD) run build
GOBUILD=$(GOCMD) build
GOFMT=gofmt
BINARY_NAME=thunderdome-planning-poker
BINARY_UNIX=$(BINARY_NAME)_unix
BINARY_WINDOWS=thunderdome-planning-poker.exe
GORELEASER=goreleaser release --rm-dist
DOCKER_TAG=stevenweathers/thunderdome-planning-poker:latest

all: build
build-deps: 
	$(NPMBUILD)
	$(STATICPACKCMD)

build: 
	$(NPMBUILD)
	$(STATICPACKCMD)
	$(GOBUILD) -o $(BINARY_NAME) -v

clean: 
	$(GOCLEAN)
	rm -f $(BINARY_NAME)
	rm -f $(BINARY_UNIX)
	rm -f $(BINARY_WINDOWS)
	rm -f *-packr.go
	rm -rf dist
	rm -rf packrd

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

run:
	DB_HOST="localhost" APP_DOMAIN="localhost" ./$(BINARY_NAME)

dev: 
	$(NPMBUILD)
	$(STATICPACKCMD)
	$(GOBUILD) -o $(BINARY_NAME) -v
	DB_HOST="localhost" APP_DOMAIN="localhost" ./$(BINARY_NAME)

release:
	$(GORELEASER)

release-dry:
	$(GORELEASER) --skip-publish

build-image:
	docker build ./ -t $(DOCKER_TAG)

push-image:
	docker push $(DOCKER_TAG)