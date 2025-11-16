FROM golang:1.23-alpine AS base
RUN apk add build-base
RUN echo $GOPATH
COPY . /app
ARG VERSION="4.13.0"



RUN set -x \
    && apk add --no-cache git \
    && git clone --branch "v${VERSION}" --depth 1 --single-branch https://github.com/golang-migrate/migrate /tmp/go-migrate

RUN apk add bash


WORKDIR /tmp/go-migrate

#ENV GOPROXY="https://goproxy.cn,direct"

RUN set -x \
    && CGO_ENABLED=0 go build -tags 'postgres' -ldflags="-s -w" -o ./migrate ./cmd/migrate \
    && ./migrate -version

RUN cp /tmp/go-migrate/migrate /usr/bin/migrate

WORKDIR /app


ENTRYPOINT [ "/bin/sh", "-c", "/usr/bin/migrate -path /app/migrations -database \"$(cat /app/.env | grep '^DATABASE_HOST=' | awk -F'=' '{print substr($0, index($0, $2))}')\" -verbose up" ]

