FROM alpine:latest
WORKDIR /app
RUN apk add --no-cache \
    ocaml build-base dune musl-dev

COPY loops.ml .
RUN ocamlopt -o loops unix.cmxa loops.ml

ENTRYPOINT ["./loops"]
