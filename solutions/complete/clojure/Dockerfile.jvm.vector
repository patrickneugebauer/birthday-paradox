FROM clojure:temurin-25-alpine
WORKDIR /app

COPY loops.clj .

ENTRYPOINT ["clj", "-M", "loops.clj"]
