FROM golang:latest as builder
WORKDIR /app
COPY . .
RUN go mod download && go mod verify
RUN go build -o /app/app
FROM alpine:latest
WORKDIR /app
RUN apk add libc6-compat
COPY --from=builder /app/app /app/app
ENTRYPOINT ["/app/app"]
