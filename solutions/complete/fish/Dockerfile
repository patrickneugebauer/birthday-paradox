FROM alpine

# MAINTAINER Maik Ellerbrock (https://github.com/ellerbrock)

# coreutils is needed for %N flag on date
RUN apk update && apk add --no-cache fish coreutils

COPY loops.fish .

ENTRYPOINT ["fish", "loops.fish"]
