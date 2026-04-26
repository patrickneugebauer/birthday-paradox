FROM erlang:alpine

WORKDIR /root

ADD loops.erl .

RUN erlc loops.erl

ENTRYPOINT ["erl", "-noshell", "-noinput", "-pa", ".", "-run", "loops", "main", "-run", "init", "stop", "--"]
