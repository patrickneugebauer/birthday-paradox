FROM crystallang/crystal:latest-alpine
WORKDIR /app

COPY loops.cr .

RUN crystal build loops.cr -o loops --release

ENTRYPOINT ["./loops"]
