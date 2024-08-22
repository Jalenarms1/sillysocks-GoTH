build_dir=./bin/sillysocks-GoTH

templ:
	@templ generate

run: build
	@$(build_dir)

build:	templ
	@go build -o $(build_dir) ./cmd/app

watch:
	@air
