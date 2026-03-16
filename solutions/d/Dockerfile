FROM dlang2/dmd-ubuntu

WORKDIR /root

ADD loops.d .

RUN dmd loops.d -od=d -of=loops -release

ENTRYPOINT ["./loops"]
