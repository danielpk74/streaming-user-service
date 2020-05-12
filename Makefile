.PHONY: build clean deploy

build:
	env GOOS=linux go build -ldflags="-s -w" -o bin/createUser src/handlers/createUser.go
	env GOOS=linux go build -ldflags="-s -w" -o bin/loginUser src/handlers/loginUser.go

clean:
	rm -rf ./bin

deploy: clean build
	sls deploy --verbose
