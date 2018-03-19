all: build

build:
	@go build -o ilber main.go

vet:
	@go vet ./...

test:
	@go test ./...

release:
	@goxc
	@rmdir debian/

deploy: release
	@scp release/0.1/ilber_*.deb do:
	@ssh do 'sudo dpkg -i ilber_*.deb'
	@ssh do 'sudo systemctl restart ilber'

issues:
	@hub issue
	@ag --ignore=Makefile -s TODO || true
	@ag --ignore=Makefile -s FIXME || true
	@ag --ignore=Makefile -s println || true

.PHONY: all build vet test release deploy issues
