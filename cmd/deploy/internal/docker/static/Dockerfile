FROM alpine:latest

RUN apk add --no-cache ca-certificates tzdata

COPY app /

ENTRYPOINT ["/app"]
