.PHONY: test golang_test nocompile_test

test: nocompile_test golang_test

golang_test:
	go test .

nocompile_test:
	go run nocompile/main.go
