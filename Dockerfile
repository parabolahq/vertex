FROM golang:latest as builder
WORKDIR /app
# Copying go.mod to reduce docker container build time (by using cache)
COPY go.mod /app/go.mod
RUN go mod download && go mod verify
COPY . .
RUN go build -o /app/app
FROM alpine:latest
WORKDIR /app
RUN apk add libc6-compat
COPY --from=builder /app/app /app/app
ENTRYPOINT ["/app/app"]
