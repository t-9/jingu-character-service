# jingu-character-service

```
docker-compose up -d --build go
cd go
docker build -t aqourstokyo/jingu-character-service .
docker push aqourstokyo/jingu-character-service:latest

cd ..
kubectl apply -f jingu-character-service-deployment.yml  
kubectl apply -f jingu-character-service-service.yml


cd nginx
docker build -t aqourstokyo/jingu-character-service-nginx .
docker push aqourstokyo/jingu-character-service-nginx:latest

cd ..
kubectl apply -f jingu-character-service-nginx-deployment.yml  
kubectl apply -f jingu-character-service-nginx-service.yml


cd mariadb
docker build -t aqourstokyo/jingu-character-service-mariadb .
docker push aqourstokyo/jingu-character-service-mariadb:latest

cd ..
kubectl create secret generic jingu-character-service-mariadb --from-literal=root-password=xxxxxxxx
kubectl apply -f jingu-character-service-mariadb-persistent-volume-claim.yml  
kubectl apply -f jingu-character-service-mariadb-deployment.yml  
kubectl apply -f jingu-character-service-mariadb-service.yml
```
