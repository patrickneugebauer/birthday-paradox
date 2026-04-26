FROM erlang:alpine
WORKDIR /app

COPY loops.erl .

RUN erlc loops.erl

ENTRYPOINT ["erl", "-noshell", "-noinput", "-pa", ".", "-run", "loops", "main", "-run", "init", "stop", "--"]
