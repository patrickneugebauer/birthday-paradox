FROM nimlang/nim

WORKDIR /root

ADD loops.nim .

RUN nim c -d:release loops.nim

ENTRYPOINT ["./loops"]
