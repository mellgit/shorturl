# AGENTS.md

Краткий набор фактов, которые помогут будущим сессиям OpenCode быстро влиться в репозиторий.

## Быстрые факты

- Сервис на Go 1.24 + Fiber. Точка входа: `main.go` -> `cmd.Up()` из `cmd/up.go`.
- В репозитории есть юнит-тесты сервисов (`internal/shortener/service_test.go`, `internal/redirect/service_test.go`). Запуск: `go test ./...`.
- Приложению на старте нужны два файла: `config.yml` и `.env`. `config.LoadEnvConfig()` вызывает `godotenv.Load()` и падает, если `.env` отсутствует.
- Миграции Goose запускаются автоматически при старте приложения (`internal/db/postgres.go`).


## Миграции

Используется Goose. Команды Makefile зависят от переменных из `.env`, потому что Makefile делает `include .env`:

```bash
make install-deps          # устанавливает goose + swag в ./bin
make migration-add name=x  # создать новую миграцию
make migration-up          # применить ожидающие миграции
make migration-down        # откатить одну миграцию
make migration-status      # показать статус миграций
```

Важно: `internal/config/config.go` задаёт `MigrationsDSN` по умолчанию с шелловскими подстановками `$(DB_HOST)`. Для `make migration-*` нужно явно задать `POSTGRES_MIGRATIONS_DSN` в `.env`.

## Swagger / кодогенерация

```bash
make swag        # swag init -g cmd/up.go, результат пишется в docs/
```

В Makefile также есть закомментированные цели `generate-proto` и `gen-sql`; они не подключены к текущему билду.

## Подводные камни

- `make <target>` упадёт без `.env` из-за `include .env`.
- Dockerfile копирует `.env` в образ на этапе сборки, но `docker-compose.yml` монтирует `.env` поверх него во время выполнения.
- В `go.mod` goose `v3.24.2`, а в `Makefile` устанавливается `v3.23.0`.
- `config.yml` конфигурирует только логирование; остальные параметры (БД, порт API, Redis) читаются из `.env`.
