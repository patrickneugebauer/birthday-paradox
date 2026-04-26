FROM debian:stable-slim

RUN apt-get update && \
    apt-get install -y --no-install-recommends gambas3 && \
    rm -rf /var/lib/apt/lists/*
# rm -rf /var/lib/apt/lists/* saves 22mb

WORKDIR /app

COPY .project .
COPY .src .

RUN gbc3 && gba3 -o loops

ENTRYPOINT ["gbr3", "./loops"]
