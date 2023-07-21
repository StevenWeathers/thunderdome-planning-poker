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

This solution will require you to pass environment variables or setup the config file, as well as setup and manage the
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
go install github.com/swaggo/swag/cmd/swag@1.7.4
npm install
```

## Build with Make

```
make build
```

### OR manual steps

### Build static assets

```
npm run build
```

### Build for current OS

```
swag init -g http/http.go -o swaggerdocs
go build
```

## Running with Watch (uses webapp dist files on OS instead of embedded)

```
npm run autobuild
make dev-go
```

## Creating SQL Migrations

First install go-migrate tool

```
go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest
```

Generate new migration files

```
migrate create -ext sql -dir db/migrations SHORT_DESCRIPTIVE_FILNAME
```

## Let the Pointing Battles begin!

Run the server and visit [http://localhost:8080](http://localhost:8080)

# Adding new Localizations

**Thunderdome** supports Locale selection on the UI (Default en-US)

Adding new locale's involves just a few steps.

1. First by copying the `ui/src/i18n/en/index.ts` into the new locale folder
   at `ui/src/i18n/{locale}/index.ts` using the two letter locale code for the directory name and translating all
   the values.
2. Second, the locale will need to be added to the locales list used by switcher component
   in ```ui/config.js``` ```locales``` object
3. Run `npm run locales` to generate the new locale types used by the build process
4. commit changes and open PR