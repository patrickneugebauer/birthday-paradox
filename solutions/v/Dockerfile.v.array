FROM thevlang/vlang:alpine
WORKDIR /app

COPY loops.v .

RUN v loops.v -o loops

ENTRYPOINT ["./loops"]
