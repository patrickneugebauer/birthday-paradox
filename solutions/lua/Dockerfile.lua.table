FROM alpine

RUN apk update && apk add --no-cache lua5.4

WORKDIR /root

COPY loops.lua .

ENTRYPOINT ["lua5.4", "loops.lua"]
