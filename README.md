# rss2discord
Post RSS feeds to Discord.

## Install
I wrote this using Go 1.16. It might work on older versions, but I'm not sure.

```
go get get.cutie.cafe/rss2discord
```

## Usage
```
Usage of rss2discord:
  -data string
        Location of a database file to write. If provided, rss2discord will "remember" the last item it sent and not post again until there's a new item at the top of the feed. This file can be shared across feeds.
  -dry
        Run a dry-run: fetch the feed and write data files if applicable, but don't post anywhere.
  -feed string
        The feed to fetch from. This can be any feed type https://github.com/mmcdole/gofeed supports.
  -hook string
        The Discord Webhook to send to (i.e. https://discord.com/api/webhooks/...)
```

Example:
```
./rss2discord -feed https://store.steampowered.com/feeds/news/app/593110 -hook https://discord.com/api/webhooks/... -data data.json
```

## License

```
Copyright (C) 2021 Alexandra Frock

This program is free software: you can redistribute it and/or modify
it under the terms of the GNU Affero General Public License as
published by the Free Software Foundation, either version 3 of the
License, or (at your option) any later version.

This program is distributed in the hope that it will be useful,
but WITHOUT ANY WARRANTY; without even the implied warranty of
MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
GNU Affero General Public License for more details.

You should have received a copy of the GNU Affero General Public License
along with this program.  If not, see <https://www.gnu.org/licenses/>.
```