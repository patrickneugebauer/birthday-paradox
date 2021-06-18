FROM daewok/lisp-devel:base

WORKDIR /root

ADD loops.cl .

ENTRYPOINT ["sbcl", "--script", "loops.cl"]
