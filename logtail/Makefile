default:	test

test:	*.go
	go test -v -race ./...

fmt:
	gofmt -w .

# Run the test in an isolated environment.
fulltest:
	docker build -t hpcloud/tail .

run:
	go run -mod vendor cmd/gotail/gotail.go -F -p -l -t /home/xshrim/test/log/aa/

build:
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -mod vendor -o gotail cmd/gotail/gotail.go
