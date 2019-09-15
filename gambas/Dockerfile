FROM mdwagner/gambas:latest

WORKDIR /root

COPY .project .
COPY .src/* .src/

RUN gbc3

ENTRYPOINT ["gbx3"]
