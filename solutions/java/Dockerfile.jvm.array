FROM eclipse-temurin:25-alpine
# this defaults to 25-jdk-alpine, should I split jdk and jre?

WORKDIR /root

COPY Loops.java .

RUN javac Loops.java

ENTRYPOINT ["java", "Loops"]
