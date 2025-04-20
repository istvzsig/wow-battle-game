build:
	@go build -o cmd/wow-battle-game

run: build
	@go run .