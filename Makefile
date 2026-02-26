init:
	go run github.com/99designs/gqlgen init
generate:
	go run github.com/99designs/gqlgen  --verbose && go run ./app/models/model_tags/model_tags.go
run:
	go run main.go
run-seed:
	go run main.go -seed
tidy:
	go mod tidy