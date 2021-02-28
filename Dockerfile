FROM golang:1.15-alpine AS builder

ARG PROJECT_PATH
ARG SRC_PATH=src
ARG BIN_NAME=app
RUN apk add --no-cache git make

WORKDIR $PROJECT_PATH

COPY go.mod go.sum ./

RUN go mod download

WORKDIR $SRC_PATH
COPY $SRC_PATH/* ./
RUN go build -o $BIN_NAME


FROM alpine:3.10

RUN apk add --no-cache --update tini ca-certificates

COPY --from=builder $PROJECT_PATH/$SRC_PATH/$BIN_NAME /$BIN_NAME

EXPOSE 14000
ENTRYPOINT ["tini", "--", "/app"]
