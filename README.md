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
| `config.allowedPointValues` | CONFIG_POINTS_ALLOWED | List of available point values for creating battles. |
| `config.defaultPointValues` | CONFIG_POINTS_DEFAULT | List of default selected points for new battles. |
| `config.show_warrior_rank` | CONFIG_SHOW_RANK     | Set to enable an icon showing the rank of a warrior during battle. |
| `config.avatar_service`    | CONFIG_AVATAR_SERVICE | Avatar service used, possible values see next paragraph |

## Avatar Service configuration

Use the name from table below to configure a service - if not set, `default` is used. Each service provides further options which then can be configured by a warrior on the profile page. Once a service is configured, drop downs with the different sprites get available. The table shows all supported services and their sprites. In all cases the same ID (`ead26688-5148-4f3c-a35d-1b0117b4f2a9`) has been used creating the avatars.

| Name |           |           |           |           |           |           |           |           |           |
| ---------- | --------- | --------- | --------- | --------- | --------- | --------- | --------- | --------- | --------- |
| `default`  |           |           |           |           |           |           |           |           |           |
|            | ![image](https://api.adorable.io/avatars/48/ead26688-5148-4f3c-a35d-1b0117b4f2a9.png) |
| `dicebear` | male | female | human | identicon | bottts | avataaars | jdenticon | gridy | code |
|            | ![image](https://avatars.dicebear.com/api/male/ead26688-5148-4f3c-a35d-1b0117b4f2a9.svg?w=48) | ![image](https://avatars.dicebear.com/api/female/ead26688-5148-4f3c-a35d-1b0117b4f2a9.svg?w=48) | ![image](https://avatars.dicebear.com/api/human/ead26688-5148-4f3c-a35d-1b0117b4f2a9.svg?w=48) | ![image](https://avatars.dicebear.com/api/identicon/ead26688-5148-4f3c-a35d-1b0117b4f2a9.svg?w=48) | ![image](https://avatars.dicebear.com/api/bottts/ead26688-5148-4f3c-a35d-1b0117b4f2a9.svg?w=48) | ![image](https://avatars.dicebear.com/api/avataaars/ead26688-5148-4f3c-a35d-1b0117b4f2a9.svg?w=48) | ![image](https://avatars.dicebear.com/api/jdenticon/ead26688-5148-4f3c-a35d-1b0117b4f2a9.svg?w=48) | ![image](https://avatars.dicebear.com/api/gridy/ead26688-5148-4f3c-a35d-1b0117b4f2a9.svg?w=48) | ![image](https://avatars.dicebear.com/api/code/ead26688-5148-4f3c-a35d-1b0117b4f2a9.svg?w=48) |
| `gravatar` | mp | identicon | monsterid | wavatar | retro | robohash | | | |
|            | ![image](https://gravatar.com/avatar/ead26688-5148-4f3c-a35d-1b0117b4f2a9?s=48&d=mp&r=g) | ![image](https://gravatar.com/avatar/ead26688-5148-4f3c-a35d-1b0117b4f2a9?s=48&d=identicon&r=g) | ![image](https://gravatar.com/avatar/ead26688-5148-4f3c-a35d-1b0117b4f2a9?s=48&d=monsterid&r=g) | ![image](https://gravatar.com/avatar/ead26688-5148-4f3c-a35d-1b0117b4f2a9?s=48&d=wavatar&r=g) | ![image](https://gravatar.com/avatar/ead26688-5148-4f3c-a35d-1b0117b4f2a9?s=48&d=retro&r=g) | ![image](https://gravatar.com/avatar/ead26688-5148-4f3c-a35d-1b0117b4f2a9?s=48&d=robohash&r=g) | | | |
| `robohash` | set1 | set2 | set3 | set4 |
|            | ![image](https://robohash.org/ead26688-5148-4f3c-a35d-1b0117b4f2a9.png?set=set1&size=48x48) | ![image](https://robohash.org/ead26688-5148-4f3c-a35d-1b0117b4f2a9.png?set=set2&size=48x48) | ![image](https://robohash.org/ead26688-5148-4f3c-a35d-1b0117b4f2a9.png?set=set3&size=48x48) | ![image](https://robohash.org/ead26688-5148-4f3c-a35d-1b0117b4f2a9.png?set=set4&size=48x48) |

# Let the Pointing Battles begin!

Run the server and visit [http://localhost:8080](http://localhost:8080)
