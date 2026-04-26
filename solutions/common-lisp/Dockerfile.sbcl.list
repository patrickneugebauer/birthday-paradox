FROM clfoundation/sbcl:2.6.2-rc-alpine3.22

WORKDIR /root

ADD loops.cl .

ENTRYPOINT ["sbcl", "--script", "loops.cl"]
