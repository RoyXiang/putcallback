# putcallback

`putcallback` is a program which handles callback from [Put.io](https://put.io/).
It will move files from Put.io to any other storage supported by [rclone](https://rclone.org/).

## Prerequisites

1. Please ensure `rclone` is installed (version >= 1.52.0)
2. Please ensure both `src` and `dest` are configured as valid rclone remotes,
   and `src` should be a remote whose type is `putio`.

## Install

Download from [Releases](https://github.com/RoyXiang/putcallback/releases/latest), or build by yourself:

```sh
env CGO_ENABLED=0 go install -trimpath -ldflags="-s -w" github.com/RoyXiang/putcallback@latest
```

## Usage

1. Setup environment variables
   * `REMOTE_SRC` and `REMOTE_DEST` (default: `src:` and `dest:`)
     * If something like `rclone move putio: drive:Downloads/` is wanted,
       set `REMOTE_SRC` to `putio:` and `REMOTE_DEST` to `drive:Downloads/`
   * `MAX_TRANSFERS` (default: the maximum number of simultaneous tasks allowed by your Put.io subscription)
   * `DELAY_BEFORE_TRANSFER` (default: `0s`)
     * Delay before moving files to the destination so that you can manipulate files on Put.io first
   * `EXCLUDE_FILETYPES` (e.g. `exe,log`)
     * Certain file types may be excluded from transferring
   * `RENAMING_STYLE`: `tv` or `anime` (default: `none`)
     * If this is set, all files downloaded by single-file transfers
       would be renamed so to be identified by media systems like Plex, Emby, Jellyfin, etc.
   * `PUSHOVER_APP_TOKEN` and `PUSHOVER_USER_TOKEN`
     * If these two are set, a notification would be sent through [Pushover](https://pushover.net/) after files transferred to `dest`.
   * `TELEGRAM_BOT_TOKEN` and `TELEGRAM_CHAT_ID`
     * If these two are set, a notification would be sent through Telegram bot after files transferred to `dest`.
2. Run the program. Set it up as a service by any means, e.g. systemd, nohup, supervisor, etc.
3. Make it accessible from outside. The program listens on `:1880` by default, set up a proxy to that port.
4. Set up callback URL on [Settings](https://app.put.io/account/settings/transfers/callback-url) page,
   e.g. `http://[your ip]:1880/`
