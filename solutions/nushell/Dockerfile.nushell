FROM alpine:latest

RUN apk add --no-cache nushell

WORKDIR /app

COPY loops.nu /app/loops.nu

ENTRYPOINT ["nu", "loops.nu"]
