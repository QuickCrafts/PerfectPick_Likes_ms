build:
	@go build -o bin/PerfectPick_Likes_ms

run: build
	@./bin/PerfectPick_Likes_ms

test:
	@go test -v ./...