FROM alpine:latest

RUN apk add --no-cache elvish

WORKDIR /app

COPY loops.elv .

ENTRYPOINT ["elvish", "loops.elv"]
