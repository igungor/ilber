all: build

build:
	@go build

vet:
	@go vet ./...

test:
	@go test ./...

deploy:
	gcloud functions deploy ilber \
		--entry-point MainHandler \
		--trigger-http \
		--region europe-west3 \
		--runtime go113 \
		--env-vars-file .env.yaml \
		--allow-unauthenticated

.PHONY: all build vet test deploy
