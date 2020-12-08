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

### **Uses WebSockets and [Svelte](https://svelte.dev/) frontend framework for a truly reactive UI experience**

![image](https://user-images.githubusercontent.com/846933/95778842-eb76ef00-0c96-11eb-99d8-af5d098c12ee.png)

# Running in production

## Use latest docker image

```
docker pull stevenweathers/thunderdome-planning-poker
```

## Use latest released binary

[![](https://img.shields.io/github/v/release/stevenweathers/thunderdome-planning-poker?include_prereleases)](https://github.com/StevenWeathers/thunderdome-planning-poker/releases/latest)

# Configuration
Thunderdome may be configured through environment variables or via a yaml file `config.yaml`
located in one of:

* `/etc/thunderdome/`
* `$HOME/.config/thunderdome/`
* Current working directory

## Required configuration items

For Thunderdome to work correctly the following configuration items are required:

| Option                     | Environment Variable | Description                                | Default Value           |
| -------------------------- | -------------------- | ------------------------------------------ | ------------------------|
| `http.domain`              | APP_DOMAIN           | The domain/base URL for this instance of Thunderdome.  Used for creating URLs in emails. | thunderdome.dev |
| `http.cookie_hashkey`      | COOKIE_HASHKEY       | Secret used to make secure cookies secure. | strongest-avenger |

### Database configuration

Thunderdome uses a Postgres database to store all data, the following configuration options exist: 

| Option                     | Environment Variable | Description                                | Default Value           |
| -------------------------- | -------------------- | ------------------------------------------ | ------------------------|
| `db.host`                  | DB_HOST              | Database host name.                        | db |
| `db.port`                  | DB_PORT              | Database port number.                      | 5432 |
| `db.user`                  | DB_USER              | Database user id.                          | thor |
| `db.pass`                  | DB_PASS              | Database user password.                    | odinson |
| `db.name`                  | DB_NAME              | Database instance name.                    | thunderdome |
| `db.sslmode`               | DB_SSLMODE           | Database SSL Mode (disable, allow, prefer, require, verify-ca, verify-full). | disable |

### SMTP (Mail) server configuration

Thunderdome sends emails for user registration related activities, the following configuration options exist:

| Option                     | Environment Variable | Description                                | Default Value           |
| -------------------------- | -------------------- | ------------------------------------------ | ------------------------|
| `smtp.host`                | SMTP_HOST            | Smtp server hostname.                      | localhost |
| `smtp.port`                | SMTP_PORT            | Smtp server port number.                   | 25 |
| `smtp.user`                | SMTP_USER            | Smtp server user.                          | |
| `smtp.pass`                | SMTP_PASS            | Smtp server password.                      | |
| `smtp.secure`              | SMTP_SECURE          | Set to authenticate with the Smtp server.  | true |
| `smtp.identity`            | SMTP_IDENTITY        | Smtp server authorization identity.  Usually unset. | |
| `smtp.sender`              | SMTP_SENDER          | From address in emails sent by Thunderdome. | no-reply@thunderdome.dev |

## Optional configuration items

| Option                     | Environment Variable | Description                                | Default Value           |
| -------------------------- | -------------------- | ------------------------------------------ | ------------------------|
| `http.port`                | PORT                 | Which port to listen for HTTP connections. | 8080 |
| `http.path_prefix`         | PATH_PREFIX          | Prefix added to all application urls for shared domain use, in format of `/{prefix}` e.g. `/thunderdome` | |
| `http.secure_cookie`       | COOKIE_SECURE        | Use secure cookies or not.                 | true |
| `http.backend_cookie_name` | BACKEND_COOKIE_NAME  | The name of the backend cookie utilized for actual auth/validation | warriorId |
| `http.frontend_cookie_name`| FRONTEND_COOKIE_NAME | The name of the cookie utilized by the UI (purely for convenience not auth) | warrior |
| `analytics.enabled`        | ANALYTICS_ENABLED    | Enable/disable google analytics.           | true |
| `analytics.id`             | ANALYTICS_ID         | Google analytics identifier.               | UA-140245309-1 |
| `config.allowedPointValues` | CONFIG_POINTS_ALLOWED | List of available point values for creating battles. | 0, 1/2, 2, 3, 5, 8, 13, 20, 40, 100, ? |
| `config.defaultPointValues` | CONFIG_POINTS_DEFAULT | List of default selected points for new battles. | 1, 2, 3, 5, 8 , 13, ? |
| `config.show_warrior_rank` | CONFIG_SHOW_RANK     | Set to enable an icon showing the rank of a warrior during battle. | false |
| `config.avatar_service`    | CONFIG_AVATAR_SERVICE | Avatar service used, possible values see next paragraph | goadorable |
| `config.toast_timeout`     | CONFIG_TOAST_TIMEOUT | Number of milliseconds before notifications are hidden. | 1000 |
| `config.allow_guests`     | CONFIG_ALLOW_GUESTS | Whether or not to allow guest (anonymous) users. | true |
| `config.allow_registration`     | CONFIG_ALLOW_REGISTRATION | Whether or not to allow user registration (outside Admin). | true |
| `config.allow_jira_import`     | CONFIG_ALLOW_JIRA_IMPORT | Whether or not to allow import plans from JIRA XML. | true |
| `config.default_locale`   | CONFIG_DEFAULT_LOCALE | The default locale (language) for the UI | en |
| `config.friendly_ui_verbs`    | CONFIG_FRIENDLY_UI_VERBS | Whether or not to use more friendly UI verbs like Users instead of Warrior, e.g. Corporate friendly | false |
| `config.allow_external_api`    | CONFIG_ALLOW_EXTERNAL_API | Whether or not to allow External API access | false |
| `auth.method`              |  AUTH_METHOD   | Choose `normal` or `ldap` as authentication method.  See separate section on LDAP configuration. | normal |

### Avatar Service configuration

Use the name from table below to configure a service - if not set, `goadorable` is used. Each service provides further options which then can be configured by a warrior on the profile page. Once a service is configured, drop downs with the different sprites become available. The table shows all supported services and their sprites. In all cases the same ID (`ead26688-5148-4f3c-a35d-1b0117b4f2a9`) has been used creating the avatars.

| Name |           |           |           |           |           |           |           |           |           |
| ---------- | --------- | --------- | --------- | --------- | --------- | --------- | --------- | --------- | --------- |
| `goadorable` (internal)  |           |           |           |           |           |           |           |           |           |
|            | ![image](https://user-images.githubusercontent.com/846933/96212071-e4283d80-0f43-11eb-9f82-ff6c105f8b0a.png) |
| `govatar` (internal) | male | female |  |  |
|            | ![image](https://user-images.githubusercontent.com/846933/96212029-ce1a7d00-0f43-11eb-9e53-8ca13ba9d4b1.png) | ![image](https://user-images.githubusercontent.com/846933/96212031-ceb31380-0f43-11eb-832b-b02c275317a5.png) |  |  |
| `dicebear` | male | female | human | identicon | bottts | avataaars | jdenticon | gridy | code |
|            | ![image](https://avatars.dicebear.com/api/male/ead26688-5148-4f3c-a35d-1b0117b4f2a9.svg?w=48) | ![image](https://avatars.dicebear.com/api/female/ead26688-5148-4f3c-a35d-1b0117b4f2a9.svg?w=48) | ![image](https://avatars.dicebear.com/api/human/ead26688-5148-4f3c-a35d-1b0117b4f2a9.svg?w=48) | ![image](https://avatars.dicebear.com/api/identicon/ead26688-5148-4f3c-a35d-1b0117b4f2a9.svg?w=48) | ![image](https://avatars.dicebear.com/api/bottts/ead26688-5148-4f3c-a35d-1b0117b4f2a9.svg?w=48) | ![image](https://avatars.dicebear.com/api/avataaars/ead26688-5148-4f3c-a35d-1b0117b4f2a9.svg?w=48) | ![image](https://avatars.dicebear.com/api/jdenticon/ead26688-5148-4f3c-a35d-1b0117b4f2a9.svg?w=48) | ![image](https://avatars.dicebear.com/api/gridy/ead26688-5148-4f3c-a35d-1b0117b4f2a9.svg?w=48) | ![image](https://avatars.dicebear.com/api/code/ead26688-5148-4f3c-a35d-1b0117b4f2a9.svg?w=48) |
| `gravatar` | mp | identicon | monsterid | wavatar | retro | robohash | | | |
|            | ![image](https://gravatar.com/avatar/ead26688-5148-4f3c-a35d-1b0117b4f2a9?s=48&d=mp&r=g) | ![image](https://gravatar.com/avatar/ead26688-5148-4f3c-a35d-1b0117b4f2a9?s=48&d=identicon&r=g) | ![image](https://gravatar.com/avatar/ead26688-5148-4f3c-a35d-1b0117b4f2a9?s=48&d=monsterid&r=g) | ![image](https://gravatar.com/avatar/ead26688-5148-4f3c-a35d-1b0117b4f2a9?s=48&d=wavatar&r=g) | ![image](https://gravatar.com/avatar/ead26688-5148-4f3c-a35d-1b0117b4f2a9?s=48&d=retro&r=g) | ![image](https://gravatar.com/avatar/ead26688-5148-4f3c-a35d-1b0117b4f2a9?s=48&d=robohash&r=g) | | | |
| `robohash` | set1 | set2 | set3 | set4 |
|            | ![image](https://robohash.org/ead26688-5148-4f3c-a35d-1b0117b4f2a9.png?set=set1&size=48x48) | ![image](https://robohash.org/ead26688-5148-4f3c-a35d-1b0117b4f2a9.png?set=set2&size=48x48) | ![image](https://robohash.org/ead26688-5148-4f3c-a35d-1b0117b4f2a9.png?set=set3&size=48x48) | ![image](https://robohash.org/ead26688-5148-4f3c-a35d-1b0117b4f2a9.png?set=set4&size=48x48) |

### LDAP Configuration

If `auth.method` is set to `ldap`, then the Create Account function is disabled and authentication
is done using LDAP.  If the LDAP server authenticates a new user successfully, the Thunderdome user 
profile is automatically generated.

The following configuration options are specific to the LDAP authentication method:

| Option                      | Environment Variable | Description                                                        |
| --------------------------- | -------------------- | ------------------------------------------------------------------ |
| `auth.ldap.url`             | AUTH_LDAP_URL        | URL to LDAP server, typically `ldap://host:port`                   |
| `auth.ldap.use_tls`         | AUTH_LDAP_USE_TLS    | Create a TLS connection after establishing the initial connection. |
| `auth.ldap.bindname`        | AUTH_LDAP_BINDNAME   | Bind name / bind DN for connecting to LDAP.  Leave empty for no authentication. |
| `auth.ldap.bindpass`        | AUTH_LDAP_BINDPASS   | Password for the bind.                                             |
| `auth.ldap.basedn`          | AUTH_LDAP_BASEDN     | Base DN for the search for the user.                               |
| `auth.ldap.filter`          | AUTH_LDAP_FILTER     | Filter for searching for the user's login id.  See below.          |
| `auth.ldap.mail_attr`       | AUTH_LDAP_MAIL_ATTR  | The LDAP property containing the user's emil address.              |
| `auth.ldap.cn_attr`         | AUTH_LDAP_CN_ATTR    | The LDAP property containing the user's name.                      |

The default `filter` is `(&(objectClass=posixAccount)(mail=%s))`.  The filter must include a `%s` that will be replaced by the user's login id.
The `mail_attr` configuration option must point to the LDAP attribute containing the user's email address.  The default is `mail`. 
The `cn_attr` configuration option must point to the LDAP attribute containing the user's full name.  The default is `cn`.

On Linux, the parameters may be tested on the command line:

```
ldapsearch -H auth.ldap.url [-Z] -x [-D auth.ldap.bindname -W] -b auth.ldap.basedn 'auth.ldap.filter' dn auth.ldap.mail auth.ldap.cn
```

The `-Z` is only used if `auth.ldap.use_tls` is set, the `-D` and `-W` parameter is only used if `auth.ldap.bindname` is set.

# Developing

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
go get
go install github.com/markbates/pkger/cmd/pkger
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

## Let the Pointing Battles begin!

Run the server and visit [http://localhost:8080](http://localhost:8080)

# Adding new Localizations
Using svelte-i18n **Thunderdome** now supports Locale selection on the UI (Default en-US)

Adding new locale's involves just a couple of steps.

1. First add the locale dictionary json files in ```frontend/public/lang/default/``` and ```frontend/public/lang/friendly/``` by copying the en.json and just changing the values of all keys
1. Second, the locale will need to be added to the locales list used by switcher component in ```frontend/config.js``` ```locales``` object