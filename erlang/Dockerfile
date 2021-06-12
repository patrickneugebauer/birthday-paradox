FROM erlang

WORKDIR /root

ADD loops.erl .

RUN erlc loops.erl

ENTRYPOINT ["erl", "-noshell", "noinput", "-pa", "./erlang", "-s", "init", "stop", "-run", "loops", "main"]
