FROM node:20-alpine3.17

COPY package.json .
COPY yarn.lock .
RUN yarn install

COPY . .

CMD ["node", "index.js"]
