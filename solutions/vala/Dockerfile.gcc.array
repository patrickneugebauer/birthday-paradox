FROM debian:stable-slim
WORKDIR /app
RUN apt update \
 && apt install -y --no-install-recommends valac gcc


COPY loops.vala .

RUN valac loops.vala -o loops

ENTRYPOINT ["./loops"]
