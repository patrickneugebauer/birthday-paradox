FROM elixir:alpine
WORKDIR /app

COPY mix.exs .
COPY lib/loops.ex lib/

RUN mix escript.build

ENTRYPOINT ["./loops"]
