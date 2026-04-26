FROM alpine:latest

RUN apk add --no-cache murex

WORKDIR /app

COPY loops.mx .

ENTRYPOINT ["murex", "loops.mx"]
