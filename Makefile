ifneq (,$(wildcard backend/.env))
    include backend/.env
    export
endif

postgres-up:
	docker compose up -d postgres

postgres-down:
	docker compose down

backend:
	cd backend && go run cmd/server/main.go 

sqlc-generate:
	cd backend && sqlc generate


migrate-up:
	cd backend && goose -dir migrations postgres "$(DATABASE_URL)" up

migrate-down:
	cd backend && goose -dir migrations postgres "$(DATABASE_URL)" down

migrate-status:
	cd backend && goose -dir migrations postgres "$(DATABASE_URL)" status
