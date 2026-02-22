# Simple Docker image using pre-built binary
FROM alpine:3.18
RUN apk add --no-cache ca-certificates
COPY server/tempconv-server-linux /usr/local/bin/tempconv-server
WORKDIR /app
EXPOSE 8080
ENV PORT=8080
CMD ["/usr/local/bin/tempconv-server"]
