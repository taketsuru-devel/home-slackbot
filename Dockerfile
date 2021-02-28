ARG SRC_PATH=src
ARG BIN_PATH="/app"

FROM golang:1.15-alpine AS builder
ARG PROJECT_PATH
ARG SRC_PATH
ARG BIN_PATH

RUN apk add --no-cache git make

WORKDIR $PROJECT_PATH

COPY go.mod go.sum ./

RUN go mod download

WORKDIR $SRC_PATH
COPY $SRC_PATH/* ./
RUN go build -o $BIN_PATH
RUN ls /

FROM alpine:3.10
ARG BIN_PATH

RUN apk add --no-cache --update tini ca-certificates

# RUN addgroup -g 1000 app && adduser -D -H -u 1000 -G app -s /bin/sh app

# COPY --from=builder --chown=app:app $BIN_PATH /$BIN_PATH
COPY --from=builder $BIN_PATH $BIN_PATH
RUN ls /
# USER app:app

EXPOSE 13000
# equivalent of $BIN_PATH
ENTRYPOINT ["tini", "--", "/app"]
