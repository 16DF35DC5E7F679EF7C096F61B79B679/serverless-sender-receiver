.PHONY: build clean deploy

build:
	env GOOS=linux go build -ldflags="-s -w" -o bin/sender sender/main.go
	env GOOS=linux go build -ldflags="-s -w" -o bin/receiver receiver/main.go

clean:
	rm -rf ./bin

deploy: clean build
	sls deploy --verbose
