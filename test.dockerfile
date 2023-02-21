FROM golang:1.9

WORKDIR /go/src/github.com/SafeRE-IT/horizon
COPY . .
ENTRYPOINT ["go", "test"]