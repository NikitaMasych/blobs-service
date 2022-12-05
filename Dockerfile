FROM golang:1.19-alpine as buildbase

RUN apk add git build-base

WORKDIR /go/src/blobs
COPY go.mod .
COPY go.sum .

RUN go mod download

COPY . .
RUN GOOS=linux go build -o /usr/local/bin/blob-service main.go


FROM alpine:3.9

COPY --from=buildbase /go/src/blobs/config.yaml /usr/local/bin/config.yaml
COPY --from=buildbase /usr/local/bin/blob-service /usr/local/bin/blob-service
RUN apk add --no-cache ca-certificates

EXPOSE 8080
ENTRYPOINT ["blob-service"]
