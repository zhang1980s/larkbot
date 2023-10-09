build:
	if [ -d case-refresh.zip ]; then rm -rf case-refresh.zip; fi;
	if [ -d bootstrap ]; then rm bootstrap ; fi;
	GOOS=linux GOARCH=arm64 go build -tags lambda.norpc -o bootstrap main.go;
	zip case-refresh.zip bootstrap