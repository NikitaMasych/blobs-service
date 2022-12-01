FROM golang:1.18-alpine as buildbase

RUN apk add git build-base

WORKDIR /go/src/blobs
COPY vendor .
COPY . .

RUN GOOS=linux go build  -o /usr/local/bin/blob-service /go/src/blobs


FROM alpine:3.9

COPY --from=buildbase /usr/local/bin/blob-service /usr/local/bin/blob-service
RUN apk add --no-cache ca-certificates

ENTRYPOINT ["blob-service"]
