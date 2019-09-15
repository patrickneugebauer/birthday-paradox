FROM haskell

RUN cabal v2-update \
 && cabal v2-install random split --lib

WORKDIR /root

COPY loops.hs .

RUN ghc loops.hs -o loops -optc-O3

ENTRYPOINT ["./loops"]
