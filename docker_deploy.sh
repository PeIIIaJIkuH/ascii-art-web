docker image build -f Dockerfile -t "ascii-art-web" .

echo "----------------------------------------------------------------"

docker images

echo "----------------------------------------------------------------"

docker container run -p 9090:8080 --detach --name ascii-art-service ascii-art-web

echo "----------------------------------------------------------------"

docker ps -a

echo "----------------------------------------------------------------"

docker exec -it ascii-art-service /bin/bash