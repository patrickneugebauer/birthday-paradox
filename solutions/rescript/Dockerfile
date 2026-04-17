FROM node:alpine

WORKDIR /app

RUN npm install rescript @rescript/core rescript-nodejs

COPY src/ ./src/
COPY rescript.json .
COPY package.json .

RUN npx rescript

ENTRYPOINT ["node", "src/Loops.res.js"]
