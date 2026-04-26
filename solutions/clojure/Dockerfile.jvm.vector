FROM clojure:temurin-25-alpine

WORKDIR /root

ADD loops.clj .

ENTRYPOINT ["clj", "-M", "loops.clj"]
