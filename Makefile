build_dir="./bin/sillysocks-GoTH.exe"

tailwind:
	@npx tailwindcss -i ./tailwind.css -o ./public/tailwind.css

templ:
	@templ generate

run: build
	@$(build_dir)

build: templ
	@go build -o $(build_dir) ./cmd/app

watch: templ
	@air --build.cmd "go build -o $(build_dir) ./cmd/app" --build.bin ".\bin\sillysocks-GoTH.exe" --build.exclude_dir="node_modules"
