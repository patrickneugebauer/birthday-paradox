FROM thevlang/vlang:alpine

WORKDIR /root

ADD loops.v .

RUN v loops.v -o loops

ENTRYPOINT ["./loops"]
