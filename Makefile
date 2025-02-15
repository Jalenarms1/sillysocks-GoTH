build_dir="./bin/sillysocks-GoTH"

run: build
	$(build_dir)

build: 
	@go build -o $(build_dir) ./cmd/sillysocks
