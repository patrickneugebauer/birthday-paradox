# with ocaml iteractive
ocaml ocaml/loops.ml

# compiled
# ----------
# compile bytecode
ocamlc -o ocaml/loops ocaml/loops.ml
# compile native
ocamlopt -o ocaml/loops ocaml/loops.ml
# run executable
ocaml/loops
