![](https://github.com/StevenWeathers/thunderdome-planning-poker/workflows/Go/badge.svg)
[![](https://img.shields.io/docker/pulls/stevenweathers/thunderdome-planning-poker.svg)](https://hub.docker.com/r/stevenweathers/thunderdome-planning-poker)
[![](https://img.shields.io/github/v/release/stevenweathers/thunderdome-planning-poker?include_prereleases)](https://github.com/StevenWeathers/thunderdome-planning-poker/releases/latest)

# Thunderdome

## Remote team collaboration with realtime agile story pointing, no cost and ad free!

![image](https://user-images.githubusercontent.com/846933/95778842-eb76ef00-0c96-11eb-99d8-af5d098c12ee.png)

## Streamline your team's agile stand-up with Team Checkins

Instead of spending time discussing what you did yesterday and what you're going to do today, focus on Blockers and
other more critical details.

![image](https://user-images.githubusercontent.com/846933/146627094-1f31a277-a454-4fd1-b707-ecb95559e9ad.png)

## Agile Sprint Retrospective

Realtime agile sprint retrospectives with grouping, voting, and action items.

![image](https://user-images.githubusercontent.com/846933/173260209-3ef3299f-f1b2-41e8-802f-17d40649c66d.png)

## Agile Feature Story Mapping

Realtime agile feature story mapping with goals, columns, stories and more!

![image](https://user-images.githubusercontent.com/846933/173260211-304a973d-4ede-494f-bb7d-b7e5c86a4e6e.png)

# Running in production

## Use latest docker image

```
docker pull stevenweathers/thunderdome-planning-poker
```

## Use latest released binary

[![](https://img.shields.io/github/v/release/stevenweathers/thunderdome-planning-poker?include_prereleases)](https://github.com/StevenWeathers/thunderdome-planning-poker/releases/latest)

# Guides

- [Configuring Thunderdome](docs/CONFIGURATION.md)
- [Contributing Guide](docs/CONTRIBUTING.md)
- [Developing Guide](docs/DEVELOPING.md) for setting up local development
- [Testing Guide](docs/TESTING.md)

# Upgrading from v1 to v2 major release

If you're currently running a 1.x.x release version of Thunderdome you will need to do the following before running a
2.x.x release version. If you're creating a fresh instance of Thunderdome you can ignore this section.

- Review the completely rewritten APIs if you're using the API feature to integrate with Thunderdome.
- Run the latest available 1.x.x release, this will run any SQL migrations that levelset the SQL schema for 2.x.x.
- Run latest available 2.x.x release, this will run any SQL migrations since 2.0.0, however will not run any 1.x.x
  migrations.
- Update any integrations using the APIs as they have been completely rewritten.

# Donations

For those who would like to donate a small amount for my efforts or monthly hosting costs of Thunderdome.dev I accept
paypal.

[![paypal](https://www.paypalobjects.com/en_US/i/btn/btn_donateCC_LG.gif)](https://paypal.me/smweathers?locale.x=en_US)