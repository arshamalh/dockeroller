# Dockeroller

Dockeroller (docker controller) is an open-source project giving you the power to control your docker daemon through Telegram, it was made for fun and personal use, but it appears to have many real-world use-cases and it is a part of ChatOps world!

> [!IMPORTANT]
> Dockeroller is made to run self-hosted, so WE WON'T ACCESS or STORE ANY OF YOUR DATA in any terms.

## How to use
First of all, you should make a bot (it's name and username doesn't matter) with [bot father](https://t.me/BotFather), copy its `token`, then you can run the CLI using:
```bash
dockeroller start --token "<YOUR-TOKEN>"
```
And now your docker daemon is accessible thorough your Telegram bot.

## Security
There is no security concerns as long as you keep your telegram account safe, and Whitelist known ids:
```bash
dockeroller start -w 22,33,44
```
For finding your ID, when you message dockeroller in UnAuthorized state, it will return your ID.

We will also add a password mechanism in the near future.

## Features
- [ ] Containers
    - [x] List
    - [x] Start
    - [x] Stop
    - [x] Live logs
    - [x] Live stats
    - [x] Remove
    - [ ] Rename
- [ ] Images
    - [ ] List
    - [ ] Remove
    - [ ] Rename

## Some Screen Shots
<img src="assets/containerslist_started_one.jpeg" alt="dockeroller containers list started one" width="500"/>
<img src="assets/containerslist_stopped_one.jpeg" alt="dockeroller containers list stopped one" width="500"/>
<img src="assets/start_command.jpeg" alt="dockeroller start command" width="500"/>
