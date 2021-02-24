
vendor:
	go mod vendor
run: vendor
	go run -mod vendor main.go
build: vendor
	# CGO_ENABLED=0 GOOS=linux go build -mod vendor -a -tags netgo -ldflags '-w' -o f5m
	CGO_ENABLED=1 GOOS=linux go build -mod vendor -a -tags netgo -ldflags '-w' -o f5m
package: build
	docker build -t ebcpaas/f5m .