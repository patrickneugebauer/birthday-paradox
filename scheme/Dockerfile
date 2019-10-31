FROM inikolaev/alpine-scheme

WORKDIR /root

COPY loops.scm .
COPY run .

RUN chmod +x run

ENTRYPOINT ["./run", "loops.scm"]
