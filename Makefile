# Go parameters
GOCMD=go
NPMCMD=npm
NPMBUILD=$(NPMCMD) run build
GOBUILD=$(GOCMD) build
SWAGGERGEN=swag init -g api/api.go -o swaggerdocs
SWAGGERDOCS=swaggerdocs
GOFMT=gofmt
GOIMPORTS=goimports
BINARY_NAME=thunderdome-planning-poker
BINARY_UNIX=$(BINARY_NAME)_unix
BINARY_WINDOWS=thunderdome-planning-poker.exe

all: build
build-deps: 
	$(NPMBUILD)
	$(SWAGGERGEN)

build: 
	$(NPMBUILD)
	$(SWAGGERGEN)
	$(GOBUILD) -o $(BINARY_NAME) -v

clean: 
	$(GOCLEAN)
	rm -f $(BINARY_NAME)
	rm -f $(BINARY_UNIX)
	rm -f $(BINARY_WINDOWS)
	rm -rf dist
	rm -rf release
	rm -rf $(SWAGGERDOCS)

format:
	$(GOFMT) -s -w */*.go
	$(GOIMPORTS) -w */*.go

testgo:
	go test `go list ./... | grep -v swaggerdocs`
# Cross compilation
build-linux:
	$(SWAGGERGEN)
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 $(GOBUILD) -o $(BINARY_UNIX) -v

build-windows:
	$(SWAGGERGEN)
	CGO_ENABLED=0 GOOS=windows GOARCH=amd64 $(GOBUILD) -o $(BINARY_WINDOWS) -v

dev: 
	$(NPMBUILD)
	$(SWAGGERGEN)
	$(GOBUILD) -o $(BINARY_NAME) -v

	SMTP_ENABLED="false" DB_HOST="localhost" APP_DOMAIN=".127.0.0.1" COOKIE_SECURE="false" ./$(BINARY_NAME) live
dev-go:
	$(SWAGGERGEN)
	$(GOBUILD) -o $(BINARY_NAME) -v

	SMTP_ENABLED="false" DB_HOST="localhost" APP_DOMAIN=".127.0.0.1" COOKIE_SECURE="false" ./$(BINARY_NAME) live
run:
	SMTP_ENABLED="false" DB_HOST="localhost" APP_DOMAIN=".127.0.0.1" COOKIE_SECURE="false" ./$(BINARY_NAME)
