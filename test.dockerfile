FROM golang:1.9

WORKDIR /go/src/github.com/saychao/horizon
COPY . .
ENTRYPOINT ["go", "test"]