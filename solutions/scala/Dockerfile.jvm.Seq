FROM sbtscala/scala-sbt:eclipse-temurin-alpine-25.0.1_8_1.12.5_3.3.7

WORKDIR /root

COPY Loops.scala .
RUN scalac Loops.scala

ENTRYPOINT [ "scala", "-J--sun-misc-unsafe-memory-access=allow", "Loops" ]
