FROM golang:1.21-alpine AS build
COPY  . /app
WORKDIR /app
RUN go build -o upload-to-minio

FROM alpine:latest
COPY --from=build /app/upload-to-minio .
EXPOSE 8080
ENTRYPOINT [ "./upload-to-minio" ]

RUN apk add --no-cache bash curl

HEALTHCHECK \
    --start-period=1s \
    --interval=1s \
    --timeout=1s \
    --retries=30 \
        CMD curl --fail -s http://localhost:8080/ping || exit 1
