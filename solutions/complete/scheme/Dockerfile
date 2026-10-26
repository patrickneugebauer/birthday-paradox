FROM debian:stable-slim
WORKDIR /app
RUN apt-get update && \
    apt-get install -y --no-install-recommends mit-scheme && \
    rm -rf /var/lib/apt/lists/*
# rm -rf /var/lib/apt/lists/* saves 22mb


COPY loops.scm .

ENTRYPOINT ["mit-scheme", "--quiet", "--load", "loops.scm", "--"]
