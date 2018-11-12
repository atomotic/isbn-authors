
FROM golang:alpine AS builder
WORKDIR /go/src/app
COPY . .
RUN apk add --no-cache git
RUN go get -d -v ./...
RUN go install -v ./...

FROM alpine:latest
RUN apk --no-cache add ca-certificates
COPY --from=builder /go/bin/app /app
ENTRYPOINT ./app
LABEL Name=isbn-authors Version=1.0.0
EXPOSE 8093
