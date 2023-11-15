.PHONY: test
test:
	go test ./...

.PHONY: up
up:
	docker compose -f docker-compose.yaml up -d --build

.PHONY: down
down:
	docker compose -f docker-compose.yaml down
