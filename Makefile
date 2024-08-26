build_dir=./bin/sillysocks-GoTH

tailwind:
	@npx tailwindcss -i ./tailwind.css -o ./public/tailwind.css --watch

templ:
	@templ generate

run: build
	@$(build_dir)

build:	templ
	@go build -o $(build_dir) ./cmd/app

watch: templ
	@air --build.cmd "go build -o $(build_dir) ./cmd/app" --build.bin "$(build_dir)" --build.exclude_dir="node_modules" --build.include_dir="views/**/*"
