.PHONY: all
all: clean dist build package release

.PHONY: clean
clean :
	rm -rf ./dist

.PHONY: dist
dist :
	mkdir ./dist

.PHONY: build
build:
	go build  -o ./dist/server ./cmd/server 


DB_CONTAINER_NAME=vehicle-server-dev
POSTGRES_USER=vehicle-server
POSTGRES_PASSWORD=secret
POSTGRES_DB=vehicle-server
DATABASE_URL=postgres://$(POSTGRES_USER):$(POSTGRES_PASSWORD)@localhost:5432/$(POSTGRES_DB)

.PHONY: dev
dev:
	go run ./cmd/server \
		-listen-address=:8080 \
		-database-url=$(DATABASE_URL)

.PHONY: dev_db
dev_db:
	docker container run \
		--detach \
		--rm \
		--name=$(DB_CONTAINER_NAME) \
		--env=POSTGRES_PASSWORD=$(POSTGRES_PASSWORD) \
		--env=POSTGRES_USER=$(POSTGRES_USER) \
		--env=POSTGRES_DB=$(POSTGRES_DB) \
		--publish 5432:5432 \
		postgis/postgis:16-3.4-alpine

.PHONY: stop_dev_db
stop_dev_db:
	docker container stop $(DB_CONTAINER_NAME)

.PHONY: unit_test integration_test
unit_test:
	go test -v -cover ./...
integration_test:
	go test -v -count=1 --tags=integration ./app


IMAGE?=	cirederf1/vehicle-server
TAG?=1.0.0
.PHONY: package
package:
	docker build -t $(IMAGE):$(TAG) .

.PHONY: release
release:
	docker push $(IMAGE):$(TAG)

.PHONY: docker_push
docker_push:
	git tag $(TAG) -a -m "Release version $(TAG)"
	git push origin $(TAG)