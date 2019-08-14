#!/bin/bash
IS_EXIST="$(git diff --name-only $TRAVIS_COMMIT_RANGE | grep postgresql)"
if [ -n "$IS_EXIST" ]
then
  # Update postgresql and bot
  echo "Start update postgresql and bot ..."
  ssh -i "travis.pem" -t $SERVER_USER@$SERVER_DNS "docker-compose -f docker-compose.yml -f docker-compose.production.yml up --force-recreate -d"
  echo "Finish update postgresql and bot"
else
  # Update only bot
  ssh -i "travis.pem" -t $SERVER_USER@$SERVER_DNS "docker-compose -f docker-compose.yml -f docker-compose.production.yml up --no-deps --force-recreate -d bot"
fi
