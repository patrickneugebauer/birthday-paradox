FROM node:lts-slim

WORKDIR /root

RUN npm i -g coffeescript

ADD loops.coffee .

RUN coffee -c loops.coffee

ENTRYPOINT ["node", "loops.js"]
