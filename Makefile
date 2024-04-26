DB_URL = $(shell awk '{if ($$0 ~ /storage/) s=1; if(s == 0) gsub(/.*/,""); if ($$0 ~ /^.$$/ && s==1) exit;} /address/ {gsub(/"/, "");print $$2}' config/config.yaml)

tidy:
	go mod tidy

compose-purge:
	docker compose down
	docker volume remove deeptruth_db_postgres_data

sqlc:
	cd config; sqlc generate

migrate-up:
	migrate -path internal/db/migration -database $(DB_URL) -verbose up

migrate-down:
	migrate -path internal/db/migration -database "$(DB_URL)" -verbose down

run:
	swag init -g cmd/server/main.go
	export CONFIG_PATH=config/config-local.yaml; go run cmd/server/main.go

test:
	go test -v cover ./...

swag:
	swag init -g cmd/server/main.go

docker-build:
	sudo docker build -t deeptruth .

compose-up:
	sudo docker compose up --build

compose-down:
	sudo docker compose down

run-db:
	sudo docker compose -f docker-compose-db.yaml up -d