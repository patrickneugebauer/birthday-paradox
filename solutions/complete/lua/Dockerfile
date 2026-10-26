FROM alpine
WORKDIR /app
RUN apk update && apk add --no-cache lua5.4


COPY loops.lua .

ENTRYPOINT ["lua5.4", "loops.lua"]
