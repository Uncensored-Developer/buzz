docker-up:
	docker compose up

start-db:
	docker compose up mysql_db

migrate-up:
	docker compose up mysql_db_go_migrate

fmt:
# run gofumpt two times because it's buggy and need 2 pass
	@find . -type f -name '*.go' ! -path '*/.cookiecutter/*' ! -path './vendor/*' | xargs gofumpt -l -w
	@find . -type f -name '*.go' ! -path '*/.cookiecutter/*' ! -path './vendor/*' | xargs gofumpt -l -w