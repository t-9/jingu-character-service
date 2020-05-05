# jingu-character-service

docker-compose up -d --build go
cd go
docker build -t aqourstokyo/jingu-character-service .
docker push aqourstokyo/jingu-character-service:latest

kubectl apply -f jingu-character-service-deployment.yml  
kubectl apply -f jingu-character-service-service.yml
