FROM golang:1.18-buster as builder

ENV TZ=America/New_York
RUN ln -snf /usr/share/zoneinfo/$TZ /etc/localtime && echo $TZ > /etc/timezone

WORKDIR /goFish
ADD . ./

RUN set -x && apt-get update && DEBIAN_FRONTEND=noninteractive apt-get install -y \
  ca-certificates && \
  rm -rf /var/lib/apt/lists/*

RUN go mod download \
  && go mod vendor \
  && go mod verify

RUN go test ./...

