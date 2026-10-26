FROM node:alpine
WORKDIR /app

RUN npm install typescript @types/node

COPY loops.ts .
COPY tsconfig.json .
RUN npx tsc

ENTRYPOINT ["node", "loops.js"]
