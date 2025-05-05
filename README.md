# short url

short url - this is a service that shortens links.

## Table of Contents
- [Docker Installation](#Docker)
- [How It Works](#Jobs)
- [Stack](#Stack)
- [Struct project](#Struct)

## <a name="Docker"></a> Docker Installation
### Run
```
make up
```
### Volumes
[Configuration file](./config.yml) `/path/config.yml:/home/app/config.yml:ro`

[Environment file](./.env) `/path/.env:/home/app/.env:ro`

### Compose

The `docker-compose.yml` file contains all the necessary databases

## <a name="Jobs"></a> How It Works

First you need to log in.

What can I do?:
- create an alias + specify the expiration date
- get a list of created aliases
- delete an alias
- update the alias
- get statistics on the number of clicks on a link
- get a generated QR code using an alias


## <a name="Stack"></a> Stack

Backend
- **Golang**
- **Fiber**
- **Validator:**
- **goose:** migrations
- **JWT**
- **swagger**
- **qrcode**

Data Base
- **PostgreSQL**
- **Redis**


**Note:** swagger documentation is available at `http://localhost:3000/swagger/index.html`

## <a name="Struct"></a> Struct project

```
shorturl/
├── cmd
│   └── up.go                          # up application
├── docs                               # swagger
│   ├── docs.go
│   ├── swagger.json
│   └── swagger.yaml
├── internal
│   ├── auth                           # auth case (JWT, registrator, login)
│   │   ├── handler.go
│   │   ├── model.go
│   │   ├── repository.go
│   │   └── service.go
│   ├── config                         # init and load config
│   │   └── config.go
│   ├── db                             # init db
│   │   ├── postgres.go
│   │   └── redis.go
│   ├── middleware                     # init jwt
│   │   └── jwt.go
│   ├── redirect                       # redirect case (redirect url)
│   │   ├── handler.go
│   │   ├── model.go
│   │   ├── repository.go
│   │   └── service.go
│   └── shortener                      # shortener case (generation, save, stats, delete, update)
│       ├── handler.go
│       ├── model.go
│       ├── repository.go
│       └── service.go
├── main.go                            # run main file
├── migrations
│   ├── 20250422162305_user_auth.sql
│   ├── 20250422193145_urls.sql
│   └── 20250426171242_init_clicks.sql 
└── pkg
    └── logger                         # init and load logger
        └── logger.go
```



