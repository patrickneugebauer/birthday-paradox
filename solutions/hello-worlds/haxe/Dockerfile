# https://haxe.org/documentation/platforms/nodejs.html
FROM node:slim

WORKDIR /app

RUN apt-get update && apt-get install -y --no-install-recommends curl ca-certificates neko && \
    mkdir -p /usr/local/lib/haxe && \
    curl -fsSL https://github.com/HaxeFoundation/haxe/releases/download/4.3.7/haxe-4.3.7-linux64.tar.gz | tar -xz -C /usr/local/lib/haxe --strip-components=1 && \
    # Symlink the binaries so they are in your PATH
    ln -s /usr/local/lib/haxe/haxe /usr/local/bin/haxe && \
    ln -s /usr/local/lib/haxe/haxelib /usr/local/bin/haxelib && \
    haxelib setup /usr/local/lib/haxe/lib && \
    haxelib install hxnodejs && \
    apt-get purge -y curl && apt-get autoremove -y && rm -rf /var/lib/apt/lists/*

COPY src .
COPY build.hxml .

RUN haxe build.hxml

ENTRYPOINT ["node", "bin/js/main.js"]
