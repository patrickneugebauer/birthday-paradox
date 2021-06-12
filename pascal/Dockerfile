FROM cmplopes/alpine-freepascal

WORKDIR /root

ADD loops.pp .

RUN fpc -O4 loops.pp

ENTRYPOINT ["./loops"]
