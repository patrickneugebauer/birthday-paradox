FROM erlang:26-alpine

RUN apk update && apk add --no-cache wget ca-certificates && \
    wget -O /tmp/gleam.tar.gz https://github.com/gleam-lang/gleam/releases/download/v1.0.0/gleam-v1.0.0-x86_64-unknown-linux-musl.tar.gz && \
    tar -xzf /tmp/gleam.tar.gz -C /usr/local/bin && \
    chmod +x /usr/local/bin/gleam && \
    apk del wget ca-certificates && \
    rm -rf /tmp/gleam.tar.gz

WORKDIR /app
COPY loops.gleam .

ENTRYPOINT ["gleam", "run"]
