DB_NAME=gtaticSite
DB_USER=hanzo
DB_PASSWORD=WVO574bJJAtr
DB_HOST=localhost
DB_PORT=5432
DB_URL=postgresql://$(DB_USER):$(DB_PASSWORD)@$(DB_HOST):$(DB_PORT)/$(DB_NAME)?sslmode=disable
postgres:
	@docker run --name pgdb -p 5432:5432 -e POSTGRES_USER=$(DB_USER) -e POSTGRES_PASSWORD=$(DB_PASSWORD) -d postgres:17.5-alpine
cdb:
	@docker exec -it pgdb createdb --username=$(DB_USER) --owner=$(DB_USER) $(DB_NAME)
ddb:
	@docker exec -it pgdb dropdb --username=$(DB_USER) $(DB_NAME)
mgup:
	migrate -path internal/db/migration -database "$(DB_URL)" -verbose up
mgdown:
	migrate -path internal/db/migration -database "$(DB_URL)" -verbose down
sqlc:
	docker run --rm -v "${CURDIR}:/src" -w /src sqlc/sqlc generate
serverDropDBCommand:
	psql -U root -> DROP DATABASE portfolioed;
serverCreateDBCommand:
	psql -U root -> CREATE DATABASE portfolioed;
templfmt:
	@templ fmt .
templgen:
	@templ generate
twc:
	tailwindcss -i internal/modules/root/css/input.css -o static/styles.css --minify --watch

.PHONY: postgres new_migration
