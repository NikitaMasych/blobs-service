FROM golang:1.19-alpine as buildbase

WORKDIR /go/src/blobs

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -mod=vendor -o /usr/local/bin/blobs-service main.go

FROM alpine:3.17.0

COPY --from=buildbase /usr/local/bin/blobs-service /usr/local/bin/blobs-service
RUN apk add --no-cache ca-certificates

CMD ["blobs-service"]