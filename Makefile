up:
	docker compose up -d
build:
	docker compose build
stop:
	docker compose stop
down:
	docker compose down
ps:
	docker compose ps
fmt:
	docker compose exec backend go fmt ./...
vet:
	docker compose exec backend go vet ./...
generate:
	docker compose exec backend go generate ./...
	docker compose exec backend go run github.com/99designs/gqlgen generate
