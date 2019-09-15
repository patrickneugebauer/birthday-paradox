FROM alpine:latest

RUN apk add --no-cache bash

WORKDIR /root

ADD loops.sh .

ENTRYPOINT ["./loops.sh"]
