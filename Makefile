TARGETNAME=github-service

.PHONY: build
build: 
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o $(TARGETNAME)
	sudo docker build .
