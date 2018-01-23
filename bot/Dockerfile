FROM golang:latest

RUN apt-get update
RUN apt-get install -y ffmpeg

RUN go get -u github.com/golang/dep/cmd/dep

ADD . /go/src/github.com/FrisovanderVeen/mobot/bot/

WORKDIR /go/src/github.com/FrisovanderVeen/mobot/bot

RUN dep ensure

ENTRYPOINT [ "go", "run", "main.go" ]