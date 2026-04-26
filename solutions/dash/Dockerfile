FROM alpine

RUN apk add --no-cache dash

WORKDIR /app

COPY loops.sh .

ENTRYPOINT ["/usr/bin/dash", "loops.sh"]
