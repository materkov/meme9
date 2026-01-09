#!/usr/bin/env bash

#cd ~/mypage

#export DOCKER_HOST=ssh://mypage-ru
cd ..
docker build --platform=linux/amd64 -t 7385cbca-brainy-vulpecula.registry.twcstorage.ru/users-service/users-service:latest users-service
docker push --platform=linux/amd64 7385cbca-brainy-vulpecula.registry.twcstorage.ru/users-service/users-service:latest

#docker tag users-service:latest 

#docker compose build users-service

cd ~/mypage

export DOCKER_HOST=ssh://mypage-ru

docker compose up -d --build users-service
