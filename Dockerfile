# horizonbuild
FROM golang:1.16.2-alpine

RUN apk add --no-cache git build-base

ARG version="dirty"

WORKDIR /go/src/github.com/SafeRE-IT/horizon
COPY . .

RUN go mod download

RUN CGO_ENABLED=0 GOOS=linux go build -ldflags "-X main.version=${version}" -o /usr/local/bin/horizon github.com/saychao/horizon/cmd/horizon