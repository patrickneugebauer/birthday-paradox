FROM alpine:latest

RUN apk add --no-cache zsh

WORKDIR /app

COPY loops.zsh .

ENTRYPOINT ["/bin/zsh", "loops.zsh"]
