FROM alpine:latest

RUN apk add go chromium

WORKDIR /go/src/url-miner
COPY . .

#RUN go install github.com/garlic0x1/url-miner@main
RUN go get -d -v ./...
RUN go build

ENTRYPOINT ./url-miner -w wordlist.txt
