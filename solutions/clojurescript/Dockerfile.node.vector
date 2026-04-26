FROM clojure:temurin-25-alpine

WORKDIR /app

RUN apk add --no-cache nodejs npm

COPY deps.edn .
COPY src ./src

RUN clj -M -m cljs.main --target node --output-to main.js -c hello-world.core

ENTRYPOINT ["node", "main.js"]
