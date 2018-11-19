all: svg test

svg:
	go run main.go > main.svg

test:
	go test
	
