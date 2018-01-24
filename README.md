# Mobot: A modular discord bot

Mobot is a modular discord bot written using [discordgo](https://github.com/bwmarrin/discordgo).

## Contents

1. Install
1. Update
1. Run
1. Plugins

## Install

* Prerequisites git, docker and docker-compose installed and a discord token and a google API key
* Run:
```bash
git clone https://github.com/FrisovanderVeen/mobot
```
* Change directory into mobot/bot
* Rename config.toml.sample to config.toml
* Edit the settings
```toml
[discord]
Token = "YOUR DISCORD TOKEN"
Prefix = "YOUR PREFIX"

[youtube]
Key = "YOUR API KEY"
```
* Change directory into mobot
* Run:
```bash
docker-compose build
```

## Update

* Change directory into mobot
* Run: 
```bash
git pull origin master
docker-compose build
```

## Run

* Change directory into mobot
* Run:
```bash
docker-compose up
```

## Plugins

Plugins are modules that can be added or removed to add or remove functionality. They should be self enclosed and removing one should not break the bot or other plugins.