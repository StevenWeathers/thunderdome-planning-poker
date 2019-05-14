# Thunderdome Planning Poker
A Planning Poker application written in Golang, stores data in Postgres db


## Building and running with docker-compose (easiest solution)

Prefered way of building and running the application with DB

```
docker-compose up --build
```

## Building

To run without docker you will need to first build, then setup the postgres DB,
and pass the user, pass, name, host, and port to the application as environment variables

```
DB_HOST=
DB_PORT=
DB_USER=
DB_PASS=
DB_NAME=
```

### Install dependencies
```
go get -d -v
go get -u github.com/gobuffalo/packr/packr
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

### bundle up static assets
```
packr
```

### Build for current OS
```
go build
```

# Let the Pointing Battles begin!

Run the server and visit [http://localhost:8080](http://localhost:8080)