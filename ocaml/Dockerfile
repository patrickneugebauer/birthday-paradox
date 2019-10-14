FROM ubuntu

WORKDIR /root

RUN apt update && apt install -y ocaml

COPY loops.ml .

RUN ocamlopt -o loops unix.cmxa loops.ml

ENTRYPOINT ["./loops"]
