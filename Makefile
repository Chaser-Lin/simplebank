mysql8:
	docker run --name mysql8 -p 3306:3306 -e MYSQL_ROOT_PASSWORD=secret -d mysql:8.0

createdb:
	docker exec -it mysql8 mysql -uroot -psecret -e "create database `simple_bank`;"

dropdb:
	docker exec -it mysql8 mysql -uroot -psecret -e "drop database `simple_bank`;"

migrateup:
	migrate -path db/migration -database "mysql://root:secret@tcp(localhost:3306)/simple_bank" -verbose up

migrateup1:
	migrate -path db/migration -database "mysql://root:secret@tcp(localhost:3306)/simple_bank" -verbose up 1

migratedown:
	migrate -path db/migration -database "mysql://root:secret@tcp(localhost:3306)/simple_bank" -verbose down

migratedown1:
	migrate -path db/migration -database "mysql://root:secret@tcp(localhost:3306)/simple_bank" -verbose down 1

sqlc:
	sqlc generate

test:
	go test -v -cover ./...

server:
	go run main.go

mock:
	mockgen -package mockdb -destination db/mock/store.go SimpleBank/db/sqlc Store

.PHONY: mysql crearedb dropdb migrateup migratedown migrateup1 migratedown1 sqlc test server mock