# Thunderdome Planning Poker
A Planning Poker application written in Golang

Goal is to build a Planning Poker application utilizing Websockets and supporting either an Embedded DB or MongoDB.

## JSON Data Model

```json
{
    "gameId": 1,
    "name": "Build a Planning Poker app",
    "creator": {
        "id": 1,
        "name": "Ricky Bobby",
    },
    "gamers": [
        {
            "id": 1,
            "name": "Sweeney Todd",
            "active": true
        }
    ],
    "stories": [
        {
            "id": 1,
            "name": "Build the Planning Poker Data Model",
            "refId": "[Story Tracking ID here]",
            "votes": [
                {
                    "voterId": 1,
                    "vote": "3"
                }
            ],
            "points": "3",
            "active": true
        }
    ],
    "active": false,
    "votingLocked": false,
}
```

## Building and running with Docker

Prefered way of building and running the application

```
docker build -t thunderdome
docker run --name thunder -p 8080:8080 thunderdome
```

## Building with Go

Will be adding Makefile soon...

```
statik -src=ui/dist -dest=vendor/thunderdome
go build
```

## Building the Webapp UI

### Project setup
```
npm install
```

### Compiles and minifies for production
```
npm run build
```

# Let the Pointing Battles begin!

Run the server and visit [http://localhost:8080](http://localhost:8080)