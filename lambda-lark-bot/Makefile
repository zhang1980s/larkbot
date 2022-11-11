build:
	if [ -d bin ]; then rm -rf bin/*; fi;
	GOARCH=amd64 GOOS=linux go build -o bin/main main.go
	zip bin/main.zip bin/main
