docker-up:
	docker compose up

start-db:
	docker compose up mysql_db

migrate-up:
	docker compose up mysql_db_go_migrate