build:
	govendor generate
	go build -o build/dp-dd-search-indexer

debug: build
	HUMAN_LOG=1 ./build/dp-dd-search-indexer

.PHONY: build debug
