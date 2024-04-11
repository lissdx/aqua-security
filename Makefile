AS_DB_NAME=as_db
DB_USER_NAME=as_2024
DB_PASSWORD=2024_as_2024
DB_URL=postgresql://$(DB_USER_NAME):$(DB_PASSWORD)@localhost:5432/$(AS_DB_NAME)?sslmode=disable

.PHONY: sqlcgen
sqlcgen:
	./as_sqlc_generator.sh


# up docker DBs
.PHONY: local_db_up
local_db_up:
	docker-compose -f docker-compose-pg.yml up -d


.PHONY: run_monitoring
run_monitoring:
	docker-compose -f as-monitoring.yml up -d


migrateupall:
	migrate -path db/migration -database $(DB_URL) -verbose up
migratedownall:
	migrate -path db/migration -database $(DB_URL) -verbose down
migrateup1:
	migrate -path db/migration -database $(DB_URL) -verbose up 1
migratedown1:
	migrate -path db/migration -database $(DB_URL) -verbose down 1