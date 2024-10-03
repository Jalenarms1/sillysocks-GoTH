build_dir="./bin/sillysocks-GoTH"

tailwind:
	@npx tailwindcss -i ./tailwind.css -o ./public/tailwind.css

run: build
	@chmod +x $(build_dir)

build: 
	@go build -o $(build_dir) ./cmd/app

watch: 
	@air --build.cmd "go build -o $(build_dir) ./cmd/app" --build.bin "./bin/sillysocks-GoTH" --build.exclude_dir="node_modules"
