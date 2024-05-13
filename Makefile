run: build
	@./bin/dreampicai

install:
	@go install

build:
	npx tailwindcss -i ./static/css/input.css -o ./static/css/output.css
	@templ generate view
	@go build -o bin/dreampicai main.go