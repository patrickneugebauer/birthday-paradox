FROM clfoundation/sbcl:2.6.2-rc-alpine3.22
WORKDIR /app

COPY loops.cl .

ENTRYPOINT ["sbcl", "--script", "loops.cl"]
