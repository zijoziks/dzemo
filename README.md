Dzemo is a Discord bot I made for my own server. It's not really the best architecture-wise and leaves much to be desired, but hey it works!

It's intended to run 24/7 via systemd. As for the features themselves, I will mostly be adding features I find necessary on my own server.

## Build

Only requirement is Go (see `go.mod` for the version).

```bash
go build -o Dzemo
```

or run it directly:

```bash
go run .
```

## Install

The intended way to run and setup the bot. This shell script will setup systemd service under a dedicated `dzemo` user, with the bot token encrypted with `systemd-creds`. You don't have to build the bot yourself, as the shell script will simply pull the latest release of the bot and all the necessary files.

```bash
curl -fsSL -o install.sh github.com/zijoziks/dzemo/raw/master/deploy/install.sh | bash
```

All other shell scripts will be located at `/opt/dzemo/deploy`.

## Update

To update the bot to latest release, run this command.

```bash
/opt/dzemo/deploy/update.sh
```

## Uninstall

To remove the bot, run this command. It will remove everything tied to Dzemo.

```bash
/opt/dzemo/deploy/uninstall.sh
```
