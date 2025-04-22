FROM golang:1.24.2-alpine AS builder

WORKDIR /etc/build-shorturl

RUN apk --no-cache add bash git make gcc gettext musl-dev

# dependencies
COPY go.mod go.sum ./
RUN go mod download

# build
COPY . .
RUN go build -o shorturl
RUN chmod +x shorturl

FROM alpine AS runner

RUN apk update && \
    apk upgrade --no-cache && \
    apk --no-cache add ffmpeg && \
    adduser -D -u 1001 -h /home/app -s /bin/sh app

WORKDIR /home/app

COPY --from=builder --chown=app:app /etc/build-shorturl/shorturl .
COPY config.yml /home/app/config.yml
COPY .env /home/app/.env
COPY ./migrations /home/app/migrations

USER app

RUN chmod +x /home/app/shorturl

CMD ["./shorturl"]
