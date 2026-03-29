FROM alpine

RUN apk update && apk add --no-cache zig
WORKDIR /app
ADD hello.zig .

RUN zig build-exe hello.zig
# ENTRYPOINT ["zig", "run", "hello.zig"]
ENTRYPOINT ["./hello"]
