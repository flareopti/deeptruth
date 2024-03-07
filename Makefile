DB_URL = $(shell awk '{if ($$0 ~ /storage/) s=1; if(s == 0) gsub(/.*/,""); if ($$0 ~ /^.$$/ && s==1) exit;} /address/ {gsub(/"/, "");print $$2}' config/config.yaml)


compose-up:
	cd config; docker compose up

compose-purge:
	cd config; docker compose down
	docker volume remove deeptruth_db_postgres_data

sqlc:
	cd config; sqlc generate

migrate-up:
	migrate -path internal/db/migration $(DB_URL) -verbose up

migrate-down:
	migrate -path internal/db/migration $(DB_URL) -verbose down

run:
	export CONFIG_PATH=config/config.yaml; go run cmd/server/main.go

test:
	go test -v cover ./...
