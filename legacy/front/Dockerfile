FROM node:12.18.4-alpine

COPY package.json .
COPY yarn.lock .
RUN yarn install

COPY . .
RUN yarn build

FROM busybox
COPY --from=0 dist /data
CMD ["cp", "-R", "/data/.", "/shared-volume"]
