FROM node:alpine

WORKDIR /root

RUN npm i -g livescript

ADD hello.ls .

RUN lsc -c hello.ls

ENTRYPOINT ["node", "hello.js"]
