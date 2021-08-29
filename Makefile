ifneq (,$(wildcard ./.env))
    include .env
    export
endif

cmd-exists-%:
	@hash $(*) > /dev/null 2>&1 || \
		(echo "ERROR: '$(*)' must be installed and available on your PATH."; exit 1)

psql: cmd-exists-psql
	psql "${DATABASE_URL}"

migrate-up:
	migrate -path ./schema -database ${DATABASE_URL} up

migrate-down:
	migrate -path ./schema -database ${DATABASE_URL} down