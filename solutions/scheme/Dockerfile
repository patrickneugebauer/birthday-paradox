FROM debian:stable-slim

RUN apt-get update && \
    apt-get install -y --no-install-recommends mit-scheme

WORKDIR /app

COPY loops.scm .

ENTRYPOINT ["mit-scheme", "--quiet", "--load", "loops.scm", "--"]
