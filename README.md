# short url
```
url-shortener/
│
├── cmd/                    # up
│
├── config/                 # configuration (env, flags, yaml)
│   └── config.go
│
├── internal/
│   ├── auth/               # auth (JWT, registrator, login)
│   │   ├── handler.go
│   │   ├── service.go
│   │   ├── repository.go
│   │   └── model.go
│   │
│   ├── shortener/          # generation, save and redirect url
│   │   ├── handler.go
│   │   ├── service.go
│   │   ├── repository.go
│   │   └── model.go
│   │
│   ├── analytics/          # processing traffic analytics
│   │   ├── handler.go
│   │   ├── service.go
│   │   ├── repository.go
│   │   └── model.go
│   │
│   ├── middleware/         # JWT, logger, rate-limiter
│   │   ├── jwt.go
│   │   ├── logger.go
│   │   └── limiter.go
│   │
│   ├── db/                 # init PostgreSQL и Redis
│   │   ├── postgres.go
│   │   └── redis.go
│   │
│   └── utils/              # sapport functions
│       └── hash.go
│
├── migrations/             # sql-migration
│
├── pkg/                    # shared libraries (tokens, validations)
│   ├── token/
│   └── validator/
│
├── Dockerfile
├── docker-compose.yml
├── .env
├── go.mod
└── README.md

```
