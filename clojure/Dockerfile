FROM clojure

RUN apt-get update \
 && apt-get install -y clojure

WORKDIR /root

ADD loops.clj .

ENTRYPOINT ["clojure", "loops.clj"]
