FROM rust

WORKDIR /root

COPY Cargo.toml .
COPY src/* src/

RUN cargo build --release

ENTRYPOINT ["./target/release/rust"]
