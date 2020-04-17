.PHONY: start
start:
	@go run main.go service.go server.go

.PHONY: test
test:
	@GO_ENV=test go test .

.PHONY: test_coverage
test_coverage:
	@GO_ENV=test go test . -cover