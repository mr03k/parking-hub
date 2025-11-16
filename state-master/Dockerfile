FROM golang:1.23-alpine AS base
RUN apk add build-base
RUN echo $GOPATH
COPY . /app

FROM base AS build
WORKDIR /app
RUN go mod download
RUN go build -o ./build ./cmd/.


FROM alpine AS prod
WORKDIR /app
COPY --from=build /app/build /app/build
RUN chmod a+x /app/build
ENTRYPOINT ["/app/build"]

