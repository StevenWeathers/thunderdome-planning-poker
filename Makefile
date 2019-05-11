# Go parameters
GOCMD=go
STATICPACKCMD=packr
NPMCMD=npm
NPMBUILD=$(NPMCMD) run build
GOBUILD=$(GOCMD) build
BINARY_NAME=thunderdome-planning-poker
BINARY_UNIX=$(BINARY_NAME)_unix
BINARY_WINDOWS=thunderdome-planning-poker.exe
GORELEASER=goreleaser release --rm-dist

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

# Cross compilation
build-linux:
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 $(GOBUILD) -o $(BINARY_UNIX) -v

build-windows:
	CGO_ENABLED=0 GOOS=windows GOARCH=amd64 $(GOBUILD) -o $(BINARY_WINDOWS) -v

run:
	./$(BINARY_NAME)

release:
	$(GORELEASER)

release-dry:
	$(GORELEASER) --skip-publish