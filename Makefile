all: build

build:
	@go build -o ilber main.go

vet:
	@go vet ./...

test:
	@go test ./...

release:
	@goxc -q -arch="amd64" -os="linux" -n="ilber" -d=release -pv=0.1
	@rmdir debian/

deploy: release
	@ansible-playbook deploy.yml

issues:
	@hub issue
	@ag --ignore=Makefile -s TODO || true
	@ag --ignore=Makefile -s FIXME || true
	@ag --ignore=Makefile -s println || true

.PHONY: all build vet test release deploy issues
