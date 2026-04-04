FROM perl:slim

RUN apt-get update && \
    apt-get install -y libtime-hires-perl && \
    rm -rf /var/lib/apt/lists/*

WORKDIR /root

COPY loops .

RUN chmod +x loops

ENTRYPOINT ["./loops"]
