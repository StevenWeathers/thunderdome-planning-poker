# Thunderdome Planning Poker
A Planning Poker application written in Golang

Goal is to build a Planning Poker application utilizing Websockets and supporting either an Embedded DB or MongoDB.

## JSON Data Model

```json
{
    "battleId": "uuid",
    "name": "Build a Planning Poker app",
    "leaderId": "uuid",
    "warriors": [
        {
            "id": "uuid",
            "name": "Sweeney Todd"
        }
    ],
    "stories": [
        {
            "id": "uuid",
            "name": "Build the Planning Poker Data Model",
            "votes": [
                {
                    "warriorId": "uuid",
                    "vote": "3"
                }
            ],
            "points": "3",
            "active": true
        }
    ],
    "votingLocked": false,
    "activePlanId": "uuid"
}
```

## Building and running with Docker

Prefered way of building and running the application

```
docker build -t thunderdome
docker run --name thunder -p 8080:8080 thunderdome
```

## Building with Go

### Install dependencies
```
go get -d -v
go get -u github.com/gobuffalo/packr/packr
npm install
```

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