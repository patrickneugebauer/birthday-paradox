FROM alpine:latest
WORKDIR /app
RUN apk add --no-cache zsh


COPY loops.zsh .

ENTRYPOINT ["/bin/zsh", "loops.zsh"]
