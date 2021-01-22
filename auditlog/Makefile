vendor:
	go mod vendor
run:
	go run -mod vendor main.go --no-header --no-body
build:
	CGO_ENABLED=0 GOOS=linux go build -a -o auditlog .
package:
	docker build -t auditlog:v1 .
deploy:
	kubectl apply -f deployment.yaml
