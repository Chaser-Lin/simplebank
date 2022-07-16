SHELL=cmd.exe
DB_URL="mysql://root:secret@tcp(localhost:3306)/simple_bank"

mysql8:
	docker run --name mysql8 -p 3306:3306 -e MYSQL_ROOT_PASSWORD=secret -d mysql:8.0

createdb:
	docker exec -it mysql8 mysql -uroot -psecret -e "create database `simple_bank`;"

dropdb:
	docker exec -it mysql8 mysql -uroot -psecret -e "drop database `simple_bank`;"

migrateup:
	migrate -path db/migration -database $(DB_URL) -verbose up

migrateup1:
	migrate -path db/migration -database $(DB_URL) -verbose up 1

migratedown:
	migrate -path db/migration -database $(DB_URL) -verbose down

migratedown1:
	migrate -path db/migration -database $(DB_URL) -verbose down 1

sqlc:
	sqlc generate

test:
	go test -v -cover ./...

server:
	go run main.go

mock:
	mockgen -package mockdb -destination db/mock/store.go SimpleBank/db/sqlc Store

proto:
	del pb\*.go
	del doc\*.json
	protoc -I=${GOPATH}/src --proto_path=proto --go_out=pb --go_opt=paths=source_relative \
	--go-grpc_out=pb --go-grpc_opt=paths=source_relative \
	--grpc-gateway_out=pb --grpc-gateway_opt=paths=source_relative \
	--openapiv2_out=doc/swagger --openapiv2_opt=allow_merge=true,merge_file_name=simple_bank \
	proto/*.proto

evans:
	evans --host localhost --port 9090 -r repl

.PHONY: mysql crearedb dropdb migrateup migratedown migrateup1 migratedown1 sqlc test server mock proto evans