build_dir="./bin/sillysocks-GoTH"

tailwind:
	@npx tailwindcss -i ./tailwind.css -o ./public/tailwind.css

run: build
	@chmod +x $(build_dir)

build: 
	@go build -o $(build_dir) ./cmd/app

watch: build
	@go run github.com/cosmtrek/air@v1.51.0