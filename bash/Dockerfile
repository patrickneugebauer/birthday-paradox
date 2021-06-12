FROM alpine:latest

RUN apk add --no-cache bash coreutils

WORKDIR /root

COPY loops.sh .

RUN chmod +x loops.sh

ENTRYPOINT ["./loops.sh"]
