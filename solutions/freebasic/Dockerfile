FROM ubuntu:22.04

WORKDIR /app

RUN apt-get update && apt-get install -y --no-install-recommends \
    curl ca-certificates xz-utils gcc libc6-dev libtinfo5 libncurses-dev && \
    curl -L https://sourceforge.net/projects/fbc/files/FreeBASIC-1.10.1/Binaries-Linux/FreeBASIC-1.10.1-linux-x86_64.tar.xz/download | tar -xJ && \
    cd FreeBASIC* && \
    ./install.sh -i && \
    cd .. && \
    rm -rf FreeBASIC* && \
    apt-get purge -y --auto-remove curl ca-certificates xz-utils && \
    apt-get clean && \
    rm -rf /var/lib/apt/lists/*

COPY hello.bas .
RUN fbc hello.bas

ENTRYPOINT ["./hello"]
