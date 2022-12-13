# Lemmy Reply Bot

This project is a simple bot that replies to comments on Lemmy. It uses Lemmy's WebSocket to get notified of any new comments, and sees if they match any regex configured in the config file. If it finds one that does, it replies with the message corresponding to that regex.

### Features

- Multiple replies in a single bot instance
- Powerful PCRE2 regular expressions for detecting triggers
- Ability to use regex capture groups in replies
- Persistent duplicate reply prevention via a filesystem store
- Uses event-based WebSocket API, which means near-instant replies and no rate limiting

### Configuration

This repo contains a file called `lemmy-reply-bot.example.toml`. This is an example config file. Copy it to `lemmy-reply-bot.toml` and edit it to fit your needs. The config contains your password, so its permissions must be set to 600 or the bot will refuse to start.

This bot uses PCRE2 for regular expressions, so you can use any of PCRE2's features, and [Regex101](https://regex101.com/) in PCRE2 mode for testing.

If any regular expressions configured in the file also match the reply messages, the bot will refuse to start because this may cause an infinite loop.

### Debugging

In order to enable debug log messages, set `LEMMY_REPLY_BOT_DEBUG=1`.

### Building and Running

First, make sure Go 1.18 or newer is installed on your system. Older versions will not work.

If you are planning to run it on the same machine as the one you're building on, simply run 

```
go build
```

And then you can start the bot with

```
./lemmy-reply-bot
```

If you want to run it on a different machine than the one you're building on, for example, a Raspberry Pi, you can build it like this:

```
GOARCH=arm64 go build
```

If your raspberry pi is 32-bit, use `arm` instead of `arm64`. Then, you can run it the same way I described above.