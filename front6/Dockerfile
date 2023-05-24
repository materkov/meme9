FROM node:18-alpine3.14

COPY package.json .
COPY yarn.lock .
RUN yarn install

COPY . .
RUN yarn build

FROM busybox
COPY --from=0 dist/ /data
CMD ["cp", "-R", "/data/.", "/apps/front"]
