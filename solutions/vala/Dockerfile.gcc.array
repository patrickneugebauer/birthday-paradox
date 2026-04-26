FROM debian:stable-slim

RUN apt update \
 && apt install -y --no-install-recommends valac gcc

WORKDIR /root

COPY loops.vala .

RUN valac loops.vala -o loops

ENTRYPOINT ["./loops"]
