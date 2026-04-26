FROM nimlang/nim
WORKDIR /app

COPY loops.nim .

RUN nim c -d:release loops.nim

ENTRYPOINT ["./loops"]
