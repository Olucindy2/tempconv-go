# Build a minimal Docker image for the tempconv-go server
FROM golang:1.21 AS builder
WORKDIR /src
COPY . .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o /tmp/tempconv-server ./server

FROM alpine:3.18
RUN apk add --no-cache ca-certificates
COPY --from=builder /tmp/tempconv-server /usr/local/bin/tempconv-server
WORKDIR /app
EXPOSE 8080
ENV PORT=8080
CMD ["/usr/local/bin/tempconv-server"]
