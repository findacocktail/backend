FROM golang:1.23.3-alpine as builder

RUN set -xe && \
    apk upgrade --update-cache --available && apk add --update gcc g++ musl-dev make && rm -rf /var/cache/apk/*

WORKDIR /home/cocktails

COPY go.mod ./
COPY go.sum ./
RUN go mod download

COPY Makefile Makefile
COPY main.go main.go
COPY internal/ internal/

RUN set -xe && go build -o bin/backend main.go

FROM alpine:latest as runner

RUN set -e; 
RUN apk upgrade --update-cache --available
RUN apk add gcc musl-dev

RUN adduser -h /home/cocktails -s /bin/bash -u 1001 -S cocktails

COPY --from=builder /home/cocktails/bin/backend /home/cocktails/backend
RUN chown -R cocktails /home/cocktails

WORKDIR /home/cocktails
USER cocktails

EXPOSE $PORT 

ENV GIN_MODE=release

CMD ["./backend"]