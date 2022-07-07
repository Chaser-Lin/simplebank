#SHELL=cmd.exe
mysql8:
	docker run --name mysql8 -p 3306:3306 -e MYSQL_ROOT_PASSWORD=secret -d mysql:8.0

createdb:
	docker exec -it mysql8 mysql -uroot -psecret -e "create database `simple_bank`;"

dropdb:
	docker exec -it mysql8 mysql -uroot -psecret -e "drop database `simple_bank`;"

migrateup:
	migrate -path db/migration -database "mysql://root:secret@tcp(localhost:3306)/simple_bank" -verbose up

migratedown:
	migrate -path db/migration -database "mysql://root:secret@tcp(localhost:3306)/simple_bank" -verbose down

sqlc:
	sqlc generate

.PHONY: mysql crearedb dropdb migrateup migrateup sqlc