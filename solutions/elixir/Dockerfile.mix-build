FROM elixir:alpine

WORKDIR /root

ADD mix.exs .
ADD lib/loops.ex lib/

RUN mix escript.build

ENTRYPOINT ["./loops"]
