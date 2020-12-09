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

deploy:
	gcloud functions deploy ilber \
		--entry-point MainHandler \
		--trigger-http \
		--region europe-west3 \
		--runtime go113 \
		--env-vars-file .env.yaml \
		--allow-unauthenticated

issues:
	@hub issue
	@ag --ignore=Makefile -s TODO || true
	@ag --ignore=Makefile -s FIXME || true
	@ag --ignore=Makefile -s println || true

.PHONY: all build vet test release deploy issues
