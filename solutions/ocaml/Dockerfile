FROM ocaml/opam:alpine

WORKDIR /home/opam

COPY loops.ml .

ENTRYPOINT ["opam", "exec", "--", "ocaml", "loops.ml"]
