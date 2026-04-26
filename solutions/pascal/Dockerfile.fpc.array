FROM freepascal/fpc:latest-slim
WORKDIR /app

COPY loops.pp .

RUN fpc -O4 loops.pp

ENTRYPOINT ["./loops"]
