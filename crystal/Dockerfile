FROM crystallang/crystal:latest

WORKDIR /root

ADD loops.cr .

RUN crystal build loops.cr -o loops --release

ENTRYPOINT ["./loops"]
