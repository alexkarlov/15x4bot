#!/bin/bash
IS_EXIST="$(git diff --name-only $TRAVIS_COMMIT_RANGE | grep postgresql)"
if [ -n "$IS_EXIST" ]
then
  # Update postgresql, bot and run migration
  echo $TRAVIS_COMMIT_RANGE
  echo "Start update postgresql and bot ..."
  ssh -i "travis.pem" -t $SERVER_USER@$SERVER_DNS "docker-compose -f docker-compose.yml -f docker-compose.production.yml up --force-recreate -d"
  ssh -i "travis.pem" -t $SERVER_USER@$SERVER_DNS "docker run -v /home/ec2-user/postgresql/migrations:/migrations --network=15x4  migrate/migrate -path=/migrations/ -database postgres://bot:@bot-db:5432/bot?sslmode=disable up"
  echo "Finish update postgresql and bot"
else
  # Update only bot
  ssh -i "travis.pem" -t $SERVER_USER@$SERVER_DNS "docker-compose -f docker-compose.yml -f docker-compose.production.yml up --no-deps --force-recreate -d bot"
fi
