# Hetzner Cloud Snapshot Backup
[![Go Report Card](https://goreportcard.com/badge/github.com/florensie/hcloud-snapshot-backup)](https://goreportcard.com/report/github.com/florensie/hcloud-snapshot-backup)
![Build](https://github.com/florensie/hcloud-snapshot-backup/workflows/Build/badge.svg)

Hetzner Cloud backups and snapshots are stored in exactly the same way.
The backup feature's flat rate will often end up more expensive however, depending on the size of your backups and price of your server.
This tool will automatically create a snapshot of every server that doesn't have backups enabled.
It will keep 7 backups per server (configurable) and automatically delete old ones.

## Usage
1. Download the [latests release](https://github.com/florensie/hcloud-snapshot-backup/releases/latest) and extract the archive or clone and build it yourself
   - OS options are `linux`, `windows` and `darwin` (MacOS)
   - `386` and `amd64` are 32-bit and 64-bit respectively
2. Generate an API token for your Hetzner Cloud project with read and write access (Security -> API tokens -> Generate API Token)
3. Modify the following values in the `.env` file or set them as environment variables
   - `HCLOUD_TOKEN`: the API token you generated
   - `KEEP_AMOUNT`: the amount of backups to keep per server
4. Run the program manually to test if it works, check the log output and your Hetzner Cloud console
5. The program creates a backup every time it is run, so you need to schedule it to run automatically yourself (with cron for example)

## Note
- Snapshots are created for all servers that don't have backups already enabled, this is not configurable
- Snapshots are created with the `autobackup` label, this is label is used to identify which snapshots should be purged (**DO NOT ADD THIS LABEL TO OTHER IMPORTANT IMAGES**)

## To Do
- Label selector: Include/exclude servers by adding the `autobackup` label to them
- Create backups for multiple servers asynchronously
- Automatically calculate if running snapshots will actually be cheaper and enable/disable the official backup feature
