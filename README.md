# VScale backup tool

This backup tool creates new and automatically removes old backups for [Vscale](https://vscale.io) instances.

It uses Vscale API for all operations and requires API token with write access.

## Installation
To install `vscale-backup` tool use go get:
```bash
go get -u github.com/alexeykiselev/vscale-backup
```

## Usage

The utility accepts the following parameters:

```
-token
-expiration 
```

The token parameter is required.
