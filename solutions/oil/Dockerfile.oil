FROM alpine:latest

WORKDIR /app

RUN apk add --no-cache curl build-base bash libstdc++ && \
    curl -LO https://oils.pub/download/oils-for-unix-0.37.0.tar.gz && \
    tar -xzf oils-for-unix-0.37.0.tar.gz && \
    cd oils-for-unix-0.37.0 && \
    ./configure && \
    _build/oils.sh && \
    ./install && \
    cd /app && \
    rm -rf oils-for-unix-0.37.0* && \
    apk del curl build-base bash

COPY loops.ysh .

ENTRYPOINT ["ysh", "loops.ysh"]
