FROM perl:slim
WORKDIR /app
RUN apt-get update && \
    apt-get install -y libtime-hires-perl && \
    rm -rf /var/lib/apt/lists/*


COPY loops .

RUN chmod +x loops

ENTRYPOINT ["./loops"]
