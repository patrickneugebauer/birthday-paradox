FROM debian:bookworm-slim

# Odin requires the full library collection (ODIN_ROOT) with core/ and other dirs.
# Also requires clang for compilation. The -file flag is needed for single-file scripts.
RUN apt-get update && apt-get install -y --no-install-recommends \
    wget ca-certificates tar clang && \
    wget -O /tmp/odin.tar.gz https://github.com/odin-lang/Odin/releases/download/dev-2026-04/odin-linux-amd64-dev-2026-04.tar.gz && \
    tar -xzf /tmp/odin.tar.gz -C /opt && \
    mv /opt/odin-linux-amd64-* /opt/odin && \
    ln -s /opt/odin/odin /usr/local/bin/odin && \
    apt-get purge -y wget ca-certificates && apt-get autoremove -y && \
    rm -rf /tmp/odin.tar.gz /var/lib/apt/lists/*

ENV ODIN_ROOT=/opt/odin

WORKDIR /app
COPY loops.odin .

ENTRYPOINT ["odin", "run", "loops.odin", "-file"]
