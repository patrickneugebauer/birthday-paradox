FROM alpine
WORKDIR /app
RUN apk update && apk add --no-cache zig
COPY hello.zig .

RUN zig build-exe hello.zig
# ENTRYPOINT ["zig", "run", "hello.zig"]
ENTRYPOINT ["./hello"]
