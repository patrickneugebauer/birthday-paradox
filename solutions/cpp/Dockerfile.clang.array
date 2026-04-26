FROM alpine
WORKDIR /app

# musl-dev is headers
RUN apk add --no-cache clang musl-dev

COPY loops.cpp .

RUN clang++ loops.cpp -O3 -o loops.out

ENTRYPOINT ["./loops.out"]
