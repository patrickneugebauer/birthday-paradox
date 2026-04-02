FROM alpine:latest

RUN apk add --no-cache \
    ocaml build-base dune musl-dev

WORKDIR /app
COPY loops.ml .
RUN ocamlopt -o loops unix.cmxa loops.ml

ENTRYPOINT ["./loops"]
