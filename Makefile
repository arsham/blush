TARGET=$(shell git describe)

deps:
	@go get github.com/Masterminds/glide
	@glide install

tmpfolder:
	mkdir -p deploy

linux: tmpfolder
	GOOS=linux GOARCH=amd64 go build -o deploy/blush main.go
	cd deploy; tar -czf blush_linux_$(TARGET).tar.gz blush ; rm blush

darwin: tmpfolder
	GOOS=darwin GOARCH=amd64 go build -o deploy/blush main.go
	cd deploy; tar -czf blush_darwin_$(TARGET).tar.gz blush ; rm blush

windows: tmpfolder
	GOOS=windows GOARCH=amd64 go build -o deploy/blush.exe main.go
	cd deploy; zip -r blush_windows_$(TARGET).zip blush.exe ; rm blush.exe

release: deps linux darwin windows

clean:
	rm -rf deploy

install: deps
	go install

update: deps
	git pull origin master

.PHONY: release linux darwin windows tmpfolder clean install deps update
