all: build

build:
	@`which go` build -v -o bin/images

release:
	@goxc -arch="amd64" -os="linux" -n="ilberbot" -d release -pv 0.1
	@rmdir debian/

deploy: release
	@ansible-playbook deploy.yml

clean:
	@rmdir debian/

.PHONY: all clean release
