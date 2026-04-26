FROM chapel/chapel

WORKDIR /app

ADD hello.chpl .

RUN chpl hello.chpl -o hello

ENTRYPOINT ["./hello"]
