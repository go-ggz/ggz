# ggz

An URL shortener service written in Golang

## Prepare

copy `.env.example` to `.env`

```sh
$ cp .env.example `.env`
```

open `.env` file 

[embedmd]:# (.env.example ini)
```ini
GGZ_DB_DRIVER=mysql
GGZ_DB_USERNAME=root
GGZ_DB_PASSWORD=123456
GGZ_DB_NAME=ggz
GGZ_DB_HOST=127.0.0.1:3307
GGZ_SERVER_ADDR=:8080
GGZ_SHORTEN_SERVER_ADDR=:8081
GGZ_DEBUG=true
GGZ_SERVER_HOST=http://localhost:8080
GGZ_SERVER_SHORTEN_HOST=http://localhost:8081
```
