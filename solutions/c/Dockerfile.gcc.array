FROM alpine
WORKDIR /app

# musl-dev is headers
RUN apk add --no-cache gcc musl-dev

COPY loops.c .

RUN cc loops.c -O3 -o loops.out

ENTRYPOINT ["./loops.out"]
