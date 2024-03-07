all: clean
	go build -o bin/shop cmd/main.go
	./bin/shop


clean: 
	rm -rf bin/*
