FROM freepascal/fpc:latest-slim

WORKDIR /root

ADD loops.pp .

RUN fpc -O4 loops.pp

ENTRYPOINT ["./loops"]
