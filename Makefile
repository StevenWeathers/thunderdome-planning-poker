# Go parameters
GOCMD=go
NPMCMD=npm
GOMODS=$(GOCMD) mod download
NPMMODS=cd ui && npm ci && cd ..
NPMBUILD=$(NPMCMD) run build --prefix ui
NPM_FORMAT=$(NPMCMD) run format --prefix ui
E2E_FORMAT=$(NPMCMD) run format --prefix e2e
GENI8N=$(NPMCMD) run locales --prefix ui
GOBUILD=$(GOCMD) build
SWAGGERDOCS=docs/swagger
SWAGGERGEN=swag init -g internal/http/http.go -o $(SWAGGERDOCS)
SWAGFORMAT=swag fmt
GOIMPORTS=goimports
BINARY_NAME=thunderdome-planning-poker
BINARY_UNIX=$(BINARY_NAME)_unix
BINARY_WINDOWS=thunderdome-planning-poker.exe

all: build
install:
	$(GOMODS)
	$(GOCMD) install github.com/swaggo/swag/cmd/swag@1.8.3
	$(NPMMODS)

build-deps:
	$(NPMBUILD)
	$(SWAGGERGEN)

build:
	$(NPMBUILD)
	$(SWAGGERGEN)
	$(GOBUILD) -o $(BINARY_NAME) -v

clean:
	rm -f $(BINARY_NAME)
	rm -f $(BINARY_UNIX)
	rm -f $(BINARY_WINDOWS)
	rm -rf ui/dist
	rm -rf release
	rm -rf $(SWAGGERDOCS)

format:
	$(GOIMPORTS) -w .
	$(SWAGFMT)
	$(NPM_FORMAT)
	$(E2E_FORMAT)

generate:
	$(GENI8N)
	$(SWAGGERGEN)

testgo:
	go test `go list ./... | grep -v $(SWAGGERDOCS)`

# Cross compilation
build-linux:
	$(SWAGGERGEN)
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 $(GOBUILD) -o $(BINARY_UNIX) -v

build-windows:
	$(SWAGGERGEN)
	CGO_ENABLED=0 GOOS=windows GOARCH=amd64 $(GOBUILD) -o $(BINARY_WINDOWS) -v

dev:
	$(GOIMPORTS) -w .
	$(SWAGFMT)
	$(GENI8N)
	$(NPM_FORMAT)
	$(NPMBUILD)
	$(SWAGGERGEN)
	$(GOBUILD) -o $(BINARY_NAME) -v

	HTTP_SECURE_PROTOCOL="false" SMTP_ENABLED="false" DB_HOST="localhost" APP_DOMAIN="localhost" COOKIE_SECURE="false" ./$(BINARY_NAME) serve --live

dev-go:
	$(GOIMPORTS) -w .
	$(SWAGFMT)
	$(SWAGGERGEN)
	$(GOBUILD) -o $(BINARY_NAME) -v

	HTTP_SECURE_PROTOCOL="false" SMTP_ENABLED="false" DB_HOST="localhost" APP_DOMAIN="localhost" COOKIE_SECURE="false" ./$(BINARY_NAME) serve --live

run:
	HTTP_SECURE_PROTOCOL="false" SMTP_ENABLED="false" DB_HOST="localhost" APP_DOMAIN="localhost" COOKIE_SECURE="false" ./$(BINARY_NAME) serve
