FROM debian:stable-slim

WORKDIR /app

RUN apt-get update && apt-get install -y curl tar && rm -rf /var/lib/apt/lists/*

RUN curl -L https://github.com/sammy-ette/Hilbish/releases/download/v2.3.4/hilbish-v2.3.4-linux-amd64.tar.gz | tar -xz \
    && mv hilbish /usr/local/bin/hilbish \
    && chmod +x /usr/local/bin/hilbish

COPY loops.lua .

ENTRYPOINT ["hilbish", "loops.lua"]
