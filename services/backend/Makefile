IMAGE_REPOSITORY := ghcr.io/izzet-mtg/storage/backend
IMAGE_VERSION := 0.0.1

.PHONY: all
all: build-docker-image

.PHONY: build-docker-image
build-docker-image:
	docker build -t "$(IMAGE_REPOSITORY):$(IMAGE_VERSION)" .
	docker tag "$(IMAGE_REPOSITORY):$(IMAGE_VERSION)" "$(IMAGE_REPOSITORY):latest"

.PHONY: push-docker-image
push-docker-image: build-docker-image
	docker push "$(IMAGE_REPOSITORY):$(IMAGE_VERSION)"

.PHONY: build-docker-compose
build-docker-compose:
	docker compose build

.PHONY: docker-compose-up
docker-compose-up:
	docker compose up

.PHONY: docker-compose-down
docker-compose-down:
	docker compose down

.PHONY: clean
clean:
	docker rmi "$(IMAGE_REPOSITORY):$(IMAGE_VERSION)" "$(IMAGE_REPOSITORY):latest"

.PHONY: generate-db-code
generate-db-code:
	sqlc generate
