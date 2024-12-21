FROM golang:1.23.3-alpine as builder

RUN set -xe && \
    apk upgrade --update-cache --available && apk add --update gcc g++ musl-dev make sqlite-dev && rm -rf /var/cache/apk/*

WORKDIR /home/cocktails

COPY go.mod ./
COPY go.sum ./
RUN go mod download

COPY Makefile Makefile
COPY main.go main.go
COPY internal/ internal/

RUN set -xe && go build -o bin/backend main.go

FROM debian:buster as runner

RUN set -e; apt-get update -y && apt-get install gcc musl-dev ca-certificates -y

RUN update-ca-certificates

RUN useradd -rm -d /home/cocktails -s /bin/bash -u 1001 -G sudo cocktails

COPY --from=builder /home/cocktails/bin/backend /home/cocktails/backend
RUN chown -R cocktails /home/cocktails

WORKDIR /home/cocktails
USER cocktails

EXPOSE $PORT 

ENV GIN_MODE=release

CMD ["./backend"]