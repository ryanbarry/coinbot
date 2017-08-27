#!/bin/sh
mkdir -p $HOME/.docker && cd $HOME/.docker
echo $SETUP_JSON | sed 's/\\n/\n/g' > setup.json
echo $CA_PEM | sed 's/\\n/\n/g' > ca.pem
echo $CERT_PEM | sed 's/\\n/\n/g' > cert.pem
echo $KEY_PEM | sed 's/\\n/\n/g' > key.pem
echo $DOCKER_AUTH > config.json
cd
docker ps --filter name=coinbot
docker rm -f coinbot
docker ps --filter name=coinbot
docker run --label com.joyent.package=g4-highcpu-128M -d --name=coinbot -e SLACK_TOKEN=$SLACK_TOKEN ryanbarry/coinbot:$CI_COMMIT_ID
docker ps --filter name=coinbot
