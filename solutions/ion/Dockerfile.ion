FROM alpine:latest

RUN apk add --no-cache cargo git libgcc && \
    cargo install --git https://github.com/redox-os/ion/ ion-shell && \
    apk del cargo git

ENV PATH="/root/.cargo/bin:${PATH}"

WORKDIR /usr/local/bin

COPY loops.ion .

ENTRYPOINT ["ion", "loops.ion"]
