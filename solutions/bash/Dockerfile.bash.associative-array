FROM alpine:latest
WORKDIR /app
RUN apk add --no-cache bash coreutils


COPY loops.sh .

RUN chmod +x loops.sh

ENTRYPOINT ["./loops.sh"]
