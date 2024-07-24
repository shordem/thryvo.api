dev:
	@nodemon --watch './**/*.go' --signal SIGTERM --exec 'go' run main.go

build:
	@go build -o thryvo

run:
	@./thryvo

build_run: build run