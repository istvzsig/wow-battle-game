build:
	go build -o cmd/wow-battle-game

start_fe:
	python3 -m http.server 5555 --directory ./web

run: build
	go run .
	start_fe