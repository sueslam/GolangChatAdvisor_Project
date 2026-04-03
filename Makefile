build-AdvisorApiFunction:
	GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -o bootstrap ./cmd/api
	cp bootstrap $(ARTIFACTS_DIR)/

test:
	go test ./...

fmt:
	go fmt ./...

sam-build: build
	sam build

sam-deploy:
	sam deploy --guided
