# GO AGI protocol implementation
#
.PHONY: test cov bench

test:
	go test -race -cover

cov:
	@go test -coverprofile=coverage.out
	@go tool cover -html=coverage.out

codecov:
	go test -race -coverprofile=coverage.txt -covermode=atomic
	bash <(curl -s https://codecov.io/bash) -t 7d4968eb-381b-4456-87bc-41bdfe331648

vet:
	@go vet -c=2

bench:
	@go test -bench=.

readme:
	mdr README.md

docmd:
	@gomarkdoc . > docs/api.md

clean:
	rm -f coverage.out
	go clean

lint:
	typos --config=.typos.toml
	golangci-lint run
