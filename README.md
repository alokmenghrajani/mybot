# mybot

`mybot` is a Slack bot I wrote for a capture the flag (CTF).

I initially forked https://github.com/rapidloop/mybot.

Given that the bot has a bunch of security bugs, I would recommend only
running it in some throwaway cloud server.

# setup

* get a Slack API token
* setup a mysql database:

      create table greeted (user varchar(50) not null primary key);

* install latex (to convert math expressions to pdf) and imagemagik (to convert pdf to png)
* install a www server (to serve static files)
* `cp config.json.sample config.json` and fill it out.
