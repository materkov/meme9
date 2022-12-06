#!/usr/bin/env bash

cd ~/mypage

docker-compose --context mypage up -d --build meme9-web
docker-compose --context mypage up -d --build meme9-front
docker-compose --context mypage up -d --build meme9-worker
