# go lang build image
FROM golang:1.22.1 AS builder
WORKDIR /app
ENV CGO_ENABLED=0 
ENV GOOS=linux
ENV GO111MODULE=on
# ENV GOPROXY=https://goproxy.cn,direct 
COPY go.mod .
COPY go.sum .
RUN go mod download -x
RUN go mod verify
COPY . .
RUN ls -laR
RUN go generate ./...
RUN go build -v  -o ./bin/ ./...


# 2nd phase
FROM alpine
WORKDIR /app
COPY --from=builder /app/bin/cmd .

COPY --from=builder /app/config.example.yaml config.example.yaml
COPY --from=builder /app/config.example.yaml config.yaml
COPY --from=builder /app/rules.yaml rules.yaml
CMD ["./cmd", "-config", "config.yaml"]




