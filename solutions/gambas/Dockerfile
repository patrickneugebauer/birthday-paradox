FROM debian:stable-slim

RUN apt-get update && \
    apt-get install -y --no-install-recommends \
    gambas3

WORKDIR /app

COPY .project .
COPY .src .

RUN gbc3 && gba3 -o loops

ENTRYPOINT ["gbr3", "./loops"]
