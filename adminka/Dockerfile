FROM node:20-alpine3.17

COPY swagger.yml .
COPY index.html .

FROM busybox
COPY --from=0 / /data
CMD ["cp", "-R", "/data/.", "/apps/adminka"]
