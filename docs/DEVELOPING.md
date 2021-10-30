# Development

## Building and running with Docker (preferred solution)

### Using Docker Compose

```
docker-compose up --build
```

### Using Docker without Compose

This solution will require you to pass environment variables or setup the config file, as well as setup and manage the DB yourself.

```
docker build ./ -f ./build/Dockerfile -t thunderdome:latest
docker run --publish 8080:8080 --name thunderdome thunderdome:latest
```

## Building

To run without docker you will need to first build, then setup the postgres DB,
and pass the user, pass, name, host, and port to the application as environment variables 
or in a config file.

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
go install github.com/swaggo/swag/cmd/swag@latest
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
swag init -g api/api.go -o swaggerdocs
go build
```

## Running with Watch (uses webapp dist files on OS instead of embedded)
```
npm run autobuild
make dev-go
```

## Let the Pointing Battles begin!

Run the server and visit [http://localhost:8080](http://localhost:8080)

# Adding new Localizations
Using svelte-i18n **Thunderdome** now supports Locale selection on the UI (Default en-US)

Adding new locale's involves just a couple of steps.

1. First add the locale dictionary json files in ```frontend/public/lang/``` and ```frontend/public/lang/default/``` and ```frontend/public/lang/friendly/``` by copying the en.json in their respective directories and just changing the values of all keys
1. Second, the locale will need to be added to the locales list used by switcher component in ```frontend/config.js``` ```locales``` object