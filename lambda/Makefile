build:
	for i in $(ls -l | grep '^d' | /usr/bin/awk '{print $9}'); do \
		if [ -d $$i/$$i.zip ]; then rm $$i/$$i.zip; fi; \
		if [ -d $$i/bootstrap ]; then rm $$i/bootstrap ; fi; \
		GOOS=linux GOARCH=arm64 go build -tags lambda.norpc -o $$i/bootstrap $$i/main.go; \
		/usr/bin/zip $$i/$$i.zip $$i/bootstrapi; \
	done
