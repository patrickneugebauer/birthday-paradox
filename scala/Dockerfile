# FROM hseeberger/scala-sbt:16.0.1_1.5.3_3.0.0
FROM hseeberger/scala-sbt:graalvm-ce-21.1.0-java11_1.5.3_3.0.0

WORKDIR /root

COPY Loops.scala .
RUN scalac Loops.scala

ENTRYPOINT [ "scala", "Loops" ]
