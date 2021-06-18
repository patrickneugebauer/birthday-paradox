FROM elixir

WORKDIR /root

ADD mix.exs .
ADD lib/ lib/

RUN mix escript.build

ENTRYPOINT ["./loops"]
