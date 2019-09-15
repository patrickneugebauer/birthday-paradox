FROM frolvlad/alpine-gcc

WORKDIR /root

ADD loops.c .

RUN cc loops.c -O3 -o loops.out

ENTRYPOINT ["./loops.out"]
