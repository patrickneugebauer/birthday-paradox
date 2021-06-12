FROM alpine

# edited from
# MAINTAINER Maik Ellerbrock (https://github.com/ellerbrock)

# coreutils is needed for %N flag on date
RUN apk update && apk add --no-cache fish coreutils

RUN sed -i -e "s/bin\/ash/usr\/bin\/fish/" /etc/passwd

ENV SHELL /usr/bin/fish

ADD loops.fish .

ENTRYPOINT ["fish", "loops.fish"]
