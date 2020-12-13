#!/bin/sh

echo 'Pulling project'

git pull --rebase

echo 'Starting mongo and redis'

docker start mongodb
docker start redis

echo 'Building project'

docker build -t oauth .

echo 'Deploying project'

docker stop login
docker run --rm --name login -d -p 8080:8080 --net test oauth

echo 'Cleanup'

docker rmi $(docker images -f "dangling=true" -q)