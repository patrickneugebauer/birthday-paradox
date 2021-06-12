FROM thevlang/vlang

WORKDIR /root

ADD loops.v .

RUN v loops.v -o loops

ENTRYPOINT ["./loops"]
