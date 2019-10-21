FROM debian

RUN apt update \
 && apt install -y valac

WORKDIR /root

COPY loops.vala .

RUN valac loops.vala -o loops

ENTRYPOINT ["./loops"]
