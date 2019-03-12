# Thunderdome Planning Poker
A Planning Poker application written in Golang

Goal is to build a Planning Poker application utilizing Websockets and supporting either an Embedded DB or MongoDB.

## JSON Data Model

```json
{
    gameId: 1,
    name: 'Build a Planning Poker app',
    creator: {
        id: 1,
        name: 'Ricky Bobby',
    },
    gamers: [
        {
            id: 1,
            name: 'Sweeney Todd',
            active: true
        }
    ],
    stories: [
        {
            id: 1,
            name: 'Build the Planning Poker Data Model'
            refId: '[Story Tracking ID here]',
            votes: [
                {
                    voterId: 1,
                    vote: '3'
                }
            ],
            points: '3',
            active: true
        }
    ],
    active: false,
    votingLocked: false,
}
```