FROM inikolaev/alpine-scheme

WORKDIR /root

COPY loops.scm .
COPY run-scheme.sh .

ENTRYPOINT ["./run-scheme.sh", "loops.scm"]
