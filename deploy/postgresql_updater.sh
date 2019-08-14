#!/bin/bash
IS_EXIST="$(git diff --name-only $TRAVIS_COMMIT_RANGE | grep postgresql)"
if [ -n "$IS_EXIST" ]
then
  echo "Start update postgresql ..."
  docker-compose -f docker-compose.yml -f docker-compose.production.yml up --force-recreate -d postgresql
  docker run -v /home/ec2-user/postgresql/migrations:/migrations --network=15x4  migrate/migrate -path=/migrations/ -database postgres://bot:@bot-db:5432/bot?sslmode=disable up
  echo "Finish update postgresql"
fi
