run: build
	@./bin/dreampicai

install:
	@go install

build:
	npx tailwindcss -i ./static/css/input.css -o ./static/css/output.css
	@templ generate view
	@go build -o bin/dreampicai main.go


up: ## Database migration up
	@go run cmd/migrate/main.go up

reset:
	@go run cmd/reset/main.go

down: ## Database migration down
	@go run cmd/migrate/main.go down

migrate: ## Migrations against the database
	@migrate create -ext sql -dir cmd/migrate/migrations $(filter-out $@,$(MAKECMDGOALS))

