# nyc-bus-transit-nearby-go-postgres-redis-leaflet

# 🚀 Mobile web app for discovering public transit in NYC. 🚀

https://github.com/coding-to-music/nyc-bus-transit-nearby-go-postgres-redis-leaflet

From / By Brian Seitz https://github.com/brnstz/bus

https://github.com/brnstz/bus

## Environment variables:

```java

```

## GitHub

```java
git init
git add .
git remote remove origin
git commit -m "first commit"
git branch -M main
git remote add origin git@github.com:coding-to-music/nyc-bus-transit-nearby-go-postgres-redis-leaflet.git
git push -u origin main
```

[![Build Status](https://travis-ci.org/brnstz/bus.svg?branch=master)](https://travis-ci.org/brnstz/bus?branch=master)

[![Token](web/src/img/token_sample_github.png)](https://token.live)

_Beta version:_ https://token.live

## Dependencies

- Go 1.6+
- PostgreSQL 9.3+ with PostGIS
- Flyway
- Redis
- NPM
- Grunt
- JQuery
- Bootstrap
- Leaflet

## Target platform

- Ubuntu 14 LTS

## Supported agencies

| Agency                 | Live departures                  |
| ---------------------- | -------------------------------- |
| MTA NYC Transit Bus    | All routes                       |
| MTA NYC Transit Subway | 123, 456, 7, ACE, BDFM, G, JZ, S |
| Staten Island Ferry    | Scheduled departures only        |

## Binaries

The full system consists of three binaries. Each binary can be configured
using environment variables and typically are run as daemons. They are
located under the `cmds/` directory.

### `busapi`

`busapi` is the queryable HTTP API. It also delivers static assets.

### `busloader`

`busloader` downloads static
[GTFS](https://developers.google.com/transit/gtfs/) files and loads those files
into the db. When it's finished loading a set of files, it updates
materialized views queried by `busapi`.

### `busprecache`

`busprecache` contacts agency-specific live data sources and writes raw
response data to Redis. `busapi` reads that data to present live
departure and vehicle location data.

## Architecture

![Token architecture](web/src/img/token_arch.png)

## Config

### Shared cache and external partner config

`busapi` and `busprecache` use these values to config Redis and external
partner sites.

| Name                       | Description                           | Default value          |
| -------------------------- | ------------------------------------- | ---------------------- |
| `BUS_REDIS_ADDR`           | `host:port` of redis                  | `localhost:6379`       |
| `BUS_REDIS_TTL`            | TTL number of seconds for Redis data  | 90                     |
| `BUS_AGENCY_IDS`           | List of agency IDs we should precache | All supported agencies |
| `BUS_MTA_BUSTIME_API_KEY`  | API key for http://bustime.mta.info/  | _None_                 |
| `BUS_MTA_DATAMINE_API_KEY` | API key for http://datamine.mta.info/ | _None_                 |

### Shared database config

All three binaries use the following database config. `busloader` must have
a writeable database, but `busprecache` and `busapi` can use a read-only replica.

| Name              | Description              | Default value    |
| ----------------- | ------------------------ | ---------------- |
| `BUS_DB_ADDR`     | `host:port` of postgres  | `localhost:5432` |
| `BUS_DB_USER`     | The username to use      | `postgres`       |
| `BUS_DB_PASSWORD` | The password to use      | empty            |
| `BUS_DB_NAME`     | The database name to use | `postgres`       |

### `busapi` config

| Name                  | Description                                         | Default value        |
| --------------------- | --------------------------------------------------- | -------------------- |
| `BUS_API_ADDR`        | The HTTP host:port we listen to                     | `0.0.0.0:8000`       |
| `BUS_WEB_DIR`         | Location of static web assets                       | `../../web/dist`     |
| `BUS_BUILD_TIMESTAMP` | Timestamp to send with static files in query string | Use API startup time |
| `BUS_LOG_TIMING`      | Log timing of certain queries                       | `false`              |

### `busloader` config

| Name               | Description                                                                             | Default value      |
| ------------------ | --------------------------------------------------------------------------------------- | ------------------ |
| `BUS_TMP_DIR`      | Path to temporary directory                                                             | `os.TempDir()`     |
| `BUS_GTFS_URLS`    | Comma-separated path to GTFS zip URLs                                                   | _None_             |
| `BUS_ROUTE_FILTER` | Comma-separated list of `route_id` values to filter on (i.e., _only_ load these routes) | _None (no filter)_ |
| `BUS_LOAD_FOREVER` | Load forever (24 hour delay between loads) if `true`, exit after first load if `false`  | `true`             |

### `busprecache` config

No specific config, just the shared cache and db configs above.

## Automation

In the `automation/` directory, there is a sample of how to fully deploy the
system. A full configuration for a deploy consists of an inventory file and a
`group_vars/` file. The included config is called `inventory_vagrant`. For
security reasons (the API keys), the vars are encrypted in this repo. You can
create your own config and deploy it locally by doing the following:

```bash

# Create vagrant server
$ cd automation/vagrant
$ vagrant up
$ cd ../..

# Overwrite group vars with defaults
$ cd automation/group_vars
$ cp defaults.yml inventory_vagrant.yml

# Add your API keys
$ vim inventory_vagrant.yml
$ cd ../..

# Deploy the system
$ cd automation
$ ./build.sh && ./deploy.sh inventory_vagrant db_install.yml db_migrations.yml api.yml web.yml loader.yml precache.yml

# If all goes well, system is available on http://localhost:8000
```

## Quickstart

Loading all data can take a long time. You can shortcut this process by
filtering for a few specific routes and data files.

```bash
# Load only the G and L train info and exit after initial load
export BUS_GTFS_URLS="http://web.mta.info/developers/data/nyct/subway/google_transit.zip"
export BUS_ROUTE_FILTER="G,L"
export BUS_LOAD_FOREVER="false"
busloader
```
