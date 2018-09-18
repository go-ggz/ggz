# ggz

[![GoDoc](https://godoc.org/github.com/go-ggz/ggz?status.svg)](https://godoc.org/github.com/go-ggz/ggz)
[![Build Status](http://drone.wu-boy.com/api/badges/go-ggz/ggz/status.svg)](http://drone.wu-boy.com/go-ggz/ggz)
[![Build status](https://ci.appveyor.com/api/projects/status/prjvsklt3io5nuhn/branch/master?svg=true)](https://ci.appveyor.com/project/appleboy/ggz/branch/master)
[![codecov](https://codecov.io/gh/go-ggz/ggz/branch/master/graph/badge.svg)](https://codecov.io/gh/go-ggz/ggz)
[![Go Report Card](https://goreportcard.com/badge/github.com/go-ggz/ggz)](https://goreportcard.com/report/github.com/go-ggz/ggz)
[![codebeat badge](https://codebeat.co/badges/0a4eff2d-c9ac-46ed-8fd7-b59942983390)](https://codebeat.co/projects/github-com-appleboy-gorush)
[![Docker Pulls](https://img.shields.io/docker/pulls/goggz/ggz.svg)](https://hub.docker.com/r/goggz/ggz/)
[![](https://images.microbadger.com/badges/image/goggz/ggz.svg)](https://microbadger.com/images/goggz/ggz "Get your own image badge on microbadger.com")
[![Release](https://github-release-version.herokuapp.com/github/go-ggz/ggz/release.svg?style=flat)](https://github.com/go-ggz/ggz/releases/latest)

An URL shortener service written in Golang.

## Features

* Support [MySQL](https://www.mysql.com/), [Postgres](https://www.postgresql.org/) or [SQLite](https://www.sqlite.org/) Database.
* Support [RESTful](https://en.wikipedia.org/wiki/Representational_state_transfer) or [GraphQL](http://graphql.org/) API.
* Support [Auth0](https://auth0.com/) Single Sign On.
* Support expose [prometheus](https://prometheus.io/) metrics.
* Support install TLS certificates from [Let's Encrypt](https://letsencrypt.org/) automatically.
* Support [QR Code](https://en.wikipedia.org/wiki/QR_code) Generator from shorten URL.
* Support local disk storage or [Minio Object Storage](https://minio.io/).
* Support linux and windows container, see [Docker Hub](https://hub.docker.com/r/goggz/ggz/tags/).

## Start app using docker-compose

See the `docker-compose.yml`

```yml
version: '3'

services:
  ggz:
    image: goggz/ggz
    restart: always
    ports:
      - 8080:8080
      - 8081:8081
    environment:
      - GGZ_DB_DRIVER=sqlite3
      - GGZ_SERVER_HOST=http://localhost:8080
      - GGZ_SERVER_SHORTEN_HOST=http://localhost:8081
      - GGZ_AUTH0_PEM_PATH=test.pem
```

## Stargazers over time

[![Stargazers over time](https://starcharts.herokuapp.com/go-ggz/ggz.svg)](https://starcharts.herokuapp.com/go-ggz/ggz)
