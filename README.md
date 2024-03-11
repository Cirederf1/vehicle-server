# Lancer la base de données:

docker container run --detach --rm --name=vehicle-server-db --env=POSTGRES_PASSWORD=secret --env=POSTGRES_USER=vehicle-server --env=POSTGRES_DB=vehicle-server --publish 5432:5432 postgis/postgis:16-3.4-alpine

# Arrêter la base de données:

docker rm -f vehicle-server-db

# Lancer le vehicle-server:

./server -listen-address=:8080 -database-url=postgres://vehicle-server:secret@localhost:5432/vehicle-server

# Créer un véhicule

```bash
curl --header "Content-Type: application/json" --data '{"latitude": 3.32,"longitude": 4.323, "shortcode":"abed", "battery": 10}' localhost:8080/vehicles | jq .
```

# Trouver les véhicules les plus proche

```bash
curl localhost:8080/vehicles\?latitude=34.2\&longitude=23.4\&limit=10
```

# Supprimer un vehicle

```bash
curl --request DELETE localhost:8080/vehicles/${VEHICLE_ID}
```
