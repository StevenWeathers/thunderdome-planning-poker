# Development

## Code Standards

### Frontend

- Follow Typescript and Svelte best practices

### Go

- [Code Review Comments guide](https://github.com/golang/go/wiki/CodeReviewComments) from the Go project should be
  followed

### Format code before commit

Run `make format` to format all go and frontend code before commit

- not doing so can result in a failed build on GitHub

## Building and running with Docker (preferred solution)

### Using Docker Compose

```
docker-compose up --build
```

### Using Docker without Compose

This solution will require you to pass environment variables or set up the config file, as well as setup and manage the
DB yourself.

```
docker build ./ -f ./build/Dockerfile -t thunderdome:latest
docker run --publish 8080:8080 --name thunderdome thunderdome:latest
```

## Building

To run without docker you will need to first build, then setup the postgres DB, and pass the user, pass, name, host, and
port to the application as environment variables or in a config file.

```
DB_HOST=
DB_PORT=
DB_USER=
DB_PASS=
DB_NAME=
```

### Install dependencies

```
go mod download
go install github.com/swaggo/swag/cmd/swag@1.8.3
cd ui && npm install
```

## Build with Make

```
make build
```

### OR manual steps

### Build static assets

```
npm run build --prefix ui
```

### Build for current OS

```
swag init -g internal/http/http.go -o docs/swagger
go build
```

## Running with Watch (uses webapp dist files on OS instead of embedded)

```
npm run autobuild --prefix ui
make dev-go
```

## Let the Pointing Battles begin!

Run the server and visit [http://localhost:8080](http://localhost:8080)

## Restful API Changes

The restful API is documented using swagger, any changes to that documentation require regenerating the docs with the
following commands and committing the updated docs with the changes.

```bash
swag fmt
swag init -g internal/http/http.go -o docs/swagger
```

## Creating SQL Migrations

First install [goose](https://github.com/pressly/goose) tool

```
go install github.com/pressly/goose/v3/cmd/goose@latest
```

Generate new migration file

```
goose -dir internal/db/migrations create SHORT_DESCRIPTIVE_FILNAME sql
```

## Adding new Localizations

**Thunderdome** supports Locale selection on the UI (Default en-US)

Adding new locale's involves just a few steps.

1. First by copying the `ui/src/i18n/en/index.ts` into the new locale folder
   at `ui/src/i18n/{locale}/index.ts` using the two letter locale code for the directory name and translating all
   the values.
2. Second, the locale will need to be added to the locales list used by switcher component
   in ```ui/config.js``` ```locales``` object
3. Run `npm run locales` in `ui` directory to generate the new locale types used by the build process
4. commit changes and open PR