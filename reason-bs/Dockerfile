FROM node:lts-slim

WORKDIR /root

ADD package.json .
ADD bsconfig.json .
ADD src src

RUN npm i
RUN npm run build

ENTRYPOINT ["node", "src/Loops.bs.js"]
