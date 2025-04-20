<!-- new migration file command -->
migrate create -ext sql -dir db/migration -seq add_users

<!-- access postgres shell -->
docker exec -it postgres17 psql -U root

<!-- github upload steps -->
git init
git add .
git commit -m "message for commit..."
git remote add origin https://github.com/your-username/your-repo-name.git
git remote -v
git branch -M main
git push -u origin main

<!-- (1) build docker image of project -->
docker build -t movie-app:latest .

<!-- (2) creating a new network so postgres and movie can have same network -->
docker network create movie-network

<!-- (3) start container from built image -->
docker run --name movie-app --network movie-network -p 8080:8080 -e GIN_MODE=release -e DB_SOURCE="postgresql://root:secret@postgres17:5432/movie-app?sslmode=disable" movie-app:latest

<!-- (4) to connect network  -->
docker network connect movie-network postgres17

<!-- to check ip address of running container -->
docker container inspect postgres17

<!-- remove a container -->
docker rm movie-app

<!-- check networks of docker containers -->
docker network ls

<!-- for more details about a docker network (bridge is name of network) -->
docker network inspect bridge

<!-- to lauch all services in a docker network at onnce using docker-compose -->

<!-- to make file executable -->
chmod +x start.sh

<!-- command to move from windows download folder to current work dir -->
mv /mnt/c/Users/Moazzan/Downloads/wait-for ./wait-for.sh

<!-- login to aws ecr using cli -->
aws ecr get-login-password | docker login --username AWS --password-stdin 861833468085.dkr.ecr.ap-south-1.amazonaws.com

<!-- sql file from dbml file -->
dbml2sql --postgres -o doc/schema.sql doc/db.dbml

<!-- before merging changes with main branch -->
git checkout -b ft/newFeature
git add .
git commit -m "new feature added"
git push origin ft/newFeature
(go to github, create pull, merge and delete feature branch)
(back in terminal)
git checkout main
git pull

<!-- using gRPC -->
make server
make evans
package pb
service movieApp
call CreateUser