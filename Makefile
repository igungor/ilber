all: build test check

build:
	@go build

vet:
	@go vet ./...

test:
	@go test -race -count=1 ./...

staticcheck:
	@staticcheck -checks inherit,-SA1019 ./...


check: vet staticcheck

deploy:
	gcloud functions deploy ilber \
		--entry-point MainHandler \
		--trigger-http \
		--region europe-west3 \
		--runtime go117 \
		--env-vars-file .env.yaml \
		--allow-unauthenticated

.PHONY: all build vet test deploy
