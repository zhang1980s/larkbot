build:
	if [ -d msg-event.zip ]; then rm -rf msg-event.zip; fi;
	if [ -d bootstrap ]; then rm bootstrap ; fi;
	GOOS=linux GOARCH=arm64 go build -tags lambda.norpc -o bootstrap main.go;
	zip msg-event.zip bootstrap