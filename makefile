user_client:
	go run cmd/user_client/main.go
user_server:
	go run cmd/user_server/main.go
ws_client:
	go run cmd/ws_client/main.go
ws_server:
	go run cmd/ws_server/main.go
gen_proto_message:
	protoc ./protos/user.proto --go_out=. --go_opt=paths=source_relative
gen_proto_grpc:
	protoc ./protos/user.proto --go-grpc_out=. --go-grpc_opt=paths=source_relative
