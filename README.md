![](https://github.com/StevenWeathers/thunderdome-planning-poker/workflows/Go/badge.svg)
![](https://github.com/StevenWeathers/thunderdome-planning-poker/workflows/Node.js%20CI/badge.svg)
![](https://github.com/StevenWeathers/thunderdome-planning-poker/workflows/Docker/badge.svg)
![](https://img.shields.io/docker/cloud/build/stevenweathers/thunderdome-planning-poker.svg)
[![](https://img.shields.io/docker/pulls/stevenweathers/thunderdome-planning-poker.svg)](https://hub.docker.com/r/stevenweathers/thunderdome-planning-poker)
[![](https://goreportcard.com/badge/github.com/stevenweathers/thunderdome-planning-poker)](https://goreportcard.com/report/github.com/stevenweathers/thunderdome-planning-poker)

# Thunderdome Planning Poker

Thunderdome is an open source agile planning poker tool in the theme of Battling for points that helps teams estimate stories.

- Planning Sessions are **Battles**
- Users are **Warriors**
- Stories are **Plans**

### **Uses WebSockets and [Svelte](https://svelte.dev/) frontend framework for a truly Reactive UI experience**

![image](https://user-images.githubusercontent.com/846933/58061038-58d62d00-7b42-11e9-9679-ebd297a51c05.png)


## Building and running with docker-compose (easiest solution)

Prefered way of building and running the application with Postgres DB

```
docker-compose up --build
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
go get
go go install github.com/markbates/pkger/cmd/pkger
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
pkger
```

### Build for current OS
```
go build
```

# Configuration
Thunderdome may be configured through environment variables or via a yaml file `config.yaml`
located in one of:

* `/etc/thunderdome/`
* `$HOME/.config/thunderdome/`
* Current working directory

The following configuration options exists:

| Option                     | Environment Variable | Description                                |
| -------------------------- | -------------------- | ------------------------------------------ |
| `http.cookie_hashkey`      | COOKIE_HASHKEY       | Secret used to make secure cookies secure. | 
| `http.port`                | PORT                 | Which port to listen for HTTP connections. |
| `http.secure_cookie`       | COOKIE_SECURE        | Use secure cookies or not.                 |
| `http.domain`              | APP_DOMAIN           | The domain/base URL for this instance of Thunderdome.  Used for creating URLs in emails. |
| `analytics.enabled`        | ANALYTICS_ENABLED    | Enable/disable google analytics.           |
| `analytics.id`             | ANALYTICS_ID         | Google analytics identifier.               |
| `db.host`                  | DB_HOST              | Database host name.                        |
| `db.port`                  | DB_PORT              | Database port number.                      |
| `db.user`                  | DB_USER              | Database user id.                          |
| `db.pass`                  | DB_PASS              | Database user password.                    |
| `db.name`                  | DB_NAME              | Database instance name.                    |
| `smtp.host`                | SMTP_HOST            | Smtp server hostname.                      |
| `smtp.port`                | SMTP_PORT            | Smtp server port number.                   |
| `smtp.secure`              | SMTP_SECURE          | Set to authenticate with the Smtp server.  |
| `smtp.identity`            | SMTP_IDENTITY        | Smtp server authorization identity.  Usually unset. |
| `smtp.sender`              | SMTP_SENDER          | From address in emails sent by Thunderdome.|
| `config.allowedPointValues` |                     | List of available point values for creating battles. |


# Let the Pointing Battles begin!

Run the server and visit [http://localhost:8080](http://localhost:8080)
