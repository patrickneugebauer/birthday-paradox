FROM clojure:lein

WORKDIR /root

ADD project.clj .
ADD src/ src/

ENTRYPOINT ["lein", "run"]
