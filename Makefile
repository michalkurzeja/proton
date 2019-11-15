default: build

.PHONY: build

build:
	@go build -o proton cmd/proton/*.go