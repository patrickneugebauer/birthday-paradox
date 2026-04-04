FROM haskell:slim

WORKDIR /root

RUN cabal v2-update \
 && cabal v2-install random split --lib

COPY loops.hs .

RUN ghc loops.hs -o loops -optc-O3 \
    -package random -package time -package split -package containers

ENTRYPOINT ["./loops"]
