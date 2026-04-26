FROM node:lts-slim
WORKDIR /app

RUN npm i -g coffeescript

COPY loops.coffee .

RUN coffee -c loops.coffee

ENTRYPOINT ["node", "loops.js"]
