FROM zenika/alpine-kotlin

WORKDIR /root

COPY loops.kt .

RUN kotlinc loops.kt -d loops.jar

ENTRYPOINT ["kotlin", "-classpath", "loops.jar", "LoopsKt"]
