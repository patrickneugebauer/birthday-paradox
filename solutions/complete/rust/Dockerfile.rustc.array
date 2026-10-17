FROM rust:slim
WORKDIR /app

COPY Cargo.toml .
COPY src/* src/

RUN cargo build --release

ENTRYPOINT ["./target/release/rust"]
