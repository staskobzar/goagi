# GO AGI protocol implementation
#
.PHONY: test cov bench

test:
	go test -race -cover

cov:
	@go test -coverprofile=coverage.out
	@go tool cover -html=coverage.out

bench:
	@go test -bench=.
