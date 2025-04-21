build:
	go build -o cmd/wow-battle-game

start_fe:
	python3 -m http.server 5555 --directory ./web

run: build
	go run . &  # Start the backend in the background
	$(MAKE) start_fe  # Start the frontend
	wait  # Wait for all background jobs to finish
