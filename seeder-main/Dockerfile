FROM golang:1.22-alpine AS base
RUN apk add build-base
RUN echo $GOPATH

COPY . /app
WORKDIR /app
RUN go build -o ./build .

ENTRYPOINT ["/app/build"]
