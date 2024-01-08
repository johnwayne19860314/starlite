ModuleName = starlite
DB_URL=postgresql://root:dummy-password@localhost:5432/starlite?sslmode=disable

network:
	docker network create bank-network

postgres:
	docker run --name postgres --network bank-network -p 5432:5432 -e POSTGRES_USER=root -e POSTGRES_PASSWORD=secret -d postgres:14-alpine

mysql:
	docker run --name mysql8 -p 3306:3306  -e MYSQL_ROOT_PASSWORD=secret -d mysql:8

createdb:
	docker exec -it postgres createdb --username=root --owner=root simple_bank

dropdb:
	docker exec -it postgres dropdb simple_bank

migrateup:
	migrate -path db/migration -database "$(DB_URL)" -verbose up

migrateup1:
	migrate -path db/migration -database "$(DB_URL)" -verbose up 1

migratedown:
	migrate -path db/migration -database "$(DB_URL)" -verbose down

migratedown1:
	migrate -path db/migration -database "$(DB_URL)" -verbose down 1

new_migration:
	migrate create -ext sql -dir db/migration -seq update_table

db_docs:
	dbdocs build doc/db.dbml

db_schema:
	dbml2sql --postgres -o doc/schema.sql doc/db.dbml

sqlc:
	sqlc generate

test:
	go test -v -cover -short ./...



mock:
	mockgen -package mockdb -destination db/mock/store.go github.com/techschool/simplebank/db/sqlc Store
	mockgen -package mockwk -destination worker/mock/distributor.go github.com/techschool/simplebank/worker TaskDistributor


evans:
	evans --host localhost --port 9090 -r repl

redis:
	docker run --name redis -p 6379:6379 -d redis:7-alpine

.PHONY: network postgres createdb dropdb migrateup migratedown migrateup1 migratedown1 new_migration db_docs db_schema sqlc test server mock proto evans redis


proto:
	scripts/proto_gen.sh

yapi:
	scripts/yapi_import.sh

openapi:
	scripts/proto_gen_yapi_import.sh
jsonYaml:
	scripts/json_to_yaml.sh

statik:
	statik -src=doc/swagger -dest=./doc