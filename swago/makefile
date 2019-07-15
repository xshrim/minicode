spec:
	swagger generate spec -mo ./swagger.json 
	#swagger validate swagger.json
	#swagger serve -F=swagger swagger.json
run:
	go run main.go
swagger:
	docker run --rm -it -e API_URL=http://192.168.43.103:9090/swagger.json -p 8081:8080 swaggerapi/swagger-ui