FROM golang:1.23-alpine AS base
RUN apk add build-base
RUN echo $GOPATH
RUN apk add git
COPY . /app


FROM base AS build

WORKDIR /app

ARG GIT_DEPLOY_TOKEN_USER
ARG GIT_DEPLOY_TOKEN_PASSWORD


ENV GO111MODULE=on
#ENV GOPROXY=https://goproxy.io,direct

ENV GOPRIVATE=git.abanppc.com/
RUN git config --global \
url."https://${GIT_DEPLOY_TOKEN_USER}:${GIT_DEPLOY_TOKEN_PASSWORD}@git.abanppc.com/".insteadOf \
"https://git.abanppc.com/"

RUN go build -o ./build ./cmd/.


FROM alpine AS prod
WORKDIR /app
COPY --from=build /app/build /app/build
RUN chmod a+x /app/build
ENTRYPOINT ["/app/build"]

