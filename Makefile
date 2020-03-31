check:
	go fmt ./...
	go vet ./...
	golangci-lint run

test:
	go test ./... -tags "unit database" -cover -race -a

sw:
	swagger generate spec -o spec.json -m
	swagger validate spec.json;

sws:
	swagger serve spec.json

comp:
	LDFLAGS="-X github.com/prometheus/common/version.Version=${CI_COMMIT_TAG:-dev} -X github.com/prometheus/common/version.Revision=${CI_COMMIT_SHA} -X github.com/prometheus/common/version.Branch=${CI_COMMIT_REF_NAME} -X github.com/prometheus/common/version.BuildUser=${GITLAB_USER_LOGIN} -X github.com/prometheus/common/version.BuildDate=$(date -u '+%Y-%m-%dT%H:%M:%S')"
	GOOS=linux CGO_ENABLED=0 go build -ldflags "all=$LDFLAGS" -a -installsuffix cgo -o app main.go

mdocker:
	docker build -f docker/Dockerfile . -t test
