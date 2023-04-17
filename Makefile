DBNAME:=integra_db
DBINST:=integra_host
MIGRATEPATH:=datastore/migrations
DB_URL:=postgres://root:password@localhost:5439/${DBNAME}?sslmode=disable

mock:
	@mockgen -source=./datastore/repository.go -destination=./rest_service/mock.go -package=rest_service
migrate_sql:
	migrate create -ext sql -dir ${MIGRATEPATH} -seq  users_schema

createdb:
	docker exec -it ${DBINST} createdb --username=root --owner=root ${DBNAME}
	#docker exec -it ${DBINST} psgl -U root ${DBNAME}

dropdb:
	docker exec -it ${DBINST} dropdb ${DBNAME}

migratedown:
	migrate -path ${MIGRATEPATH} -database "$(DB_URL)" -verbose down
migrateup:
	migrate -path ${MIGRATEPATH} -database "$(DB_URL)" -verbose up

setup:
	chmod +x create-env.sh
	./create-env.sh
rebuild:
	docker-compose up --build -d

start:
	docker-compose up -d

downV:
	docker-compose down -v

log_api:
	docker logs -f integra_api

.PHONY: mock migrate_sql createdb dropdb migrateup migratedown setup reBuild down start log_api
