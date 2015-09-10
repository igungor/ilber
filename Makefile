all: build

build:
	@go build -o ilberbot main.go

vet:
	@go vet ./...

test:
	@go test ./...

release:
	@goxc -q -arch="amd64" -os="linux" -n="ilberbot" -d release -pv 0.1
	@rmdir debian/

deploy: release
	@ansible-playbook deploy.yml

.PHONY: all build vet test release deploy
