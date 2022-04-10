FROM golang:1.18


WORKDIR /go/src/url-miner
COPY . .

RUN go install github.com/garlic0x1/url-miner@main

ENTRYPOINT url-miner -w wordlist.txt
