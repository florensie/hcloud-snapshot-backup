# Hetzner Cloud Snapshot Backup
[![Go Report Card](https://goreportcard.com/badge/github.com/florensie/hcloud-snapshot-backup)](https://goreportcard.com/report/github.com/florensie/hcloud-snapshot-backup)
![Build](https://github.com/florensie/hcloud-snapshot-backup/workflows/Build/badge.svg)

Hetzner Cloud backups and snapshots are stored in exactly the same way.
The backup feature's flat rate will often end up more expensive however, depending on the size of your backups and price of your server.
This tool will automatically create a snapshot of every server that doesn't have backups enabled.
It will keep 7 backups per server (configurable) and automatically delete old ones.

## Usage
1. Set the `HCLOUD_TOKEN` and `KEEP_AMOUNT` environment variables or use the `.env` file
2. Use cron to run the program periodically at your desired schedule

## To Do
- Create backups asynchronously
- Automatically calculate if running snapshots will actually be cheaper and enable/disable the official backup feature
