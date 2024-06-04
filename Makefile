usersvc:
	@go build -o bin/usersvc ./usersvc/cmd
	@./bin/usersvc
test:
	@go test -v ./... --count=1
build:
	@docker build -t usersvc .
run: 
	@docker run -d -p 3004:3004 usersvc:latest
proto:
	@protoc --go_out=. --go_opt=paths=source_relative --go-grpc_out=. --go-grpc_opt=paths=source_relative usersvc/proto/*.proto

.PHONY: usersvc