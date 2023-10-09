build:
	if [ -d card-event.zip ]; then rm -rf card-event.zip; fi;
	if [ -d bootstrap ]; then rm bootstrap ; fi;
	GOOS=linux GOARCH=arm64 go build -tags lambda.norpc -o bootstrap main.go;
	zip card-event.zip bootstrap