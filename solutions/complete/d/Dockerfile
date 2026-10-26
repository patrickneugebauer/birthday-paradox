FROM dlang2/dmd-ubuntu
WORKDIR /app

COPY loops.d .

# RUN dmd loops.d -od=d -of=loops -release
RUN dmd -O -release -inline -boundscheck=off loops.d -of=loops

ENTRYPOINT ["./loops"]
