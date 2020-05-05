# jingu-character-service

docker-compose up -d --build go
cd go
docker build -t aqourstokyo/jingu-app-service .
docker push aqourstokyo/jingu-app-service:latest

kubectl apply -f jingu-character-service-deployment.yml  
kubectl apply -f jingu-character-service-service.yml
