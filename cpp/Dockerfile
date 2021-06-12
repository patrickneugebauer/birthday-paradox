FROM frolvlad/alpine-gxx

WORKDIR /root

ADD loops.cpp .

RUN g++ loops.cpp -O3 -o loops.out

ENTRYPOINT ["./loops.out"]
