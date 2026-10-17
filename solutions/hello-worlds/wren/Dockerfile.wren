FROM debian:bookworm-slim

# Wren is an embedded language. wren_test executable can run scripts.
RUN apt-get update && apt-get install -y git build-essential ca-certificates && \
    git clone https://github.com/wren-lang/wren.git /tmp/wren && \
    cd /tmp/wren/projects/make && make wren_test && \
    cp /tmp/wren/bin/wren_test /usr/local/bin/wren && chmod +x /usr/local/bin/wren

WORKDIR /app
COPY loops.wren .

ENTRYPOINT ["wren", "loops.wren"]
