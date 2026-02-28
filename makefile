# Start http server
http_server:
	go run cmd/http_server/main.go
# Run sample client programm to test grpc server
user_client:
	go run cmd/user_client/main.go
# Start grpc server
user_server:
	go run cmd/user_server/main.go
# Run sample client programm to test websocket server
ws_client:
	go run cmd/ws_client/main.go
# Start websocket server
ws_server:
	go run cmd/ws_server/main.go
# Generate Go structures for protobuf messages
gen_proto_message:
	protoc ./protos/user.proto --go_out=. --go_opt=paths=source_relative
# Generate Go client and server code for protobuf service
gen_proto_grpc:
	protoc ./protos/user.proto --go-grpc_out=. --go-grpc_opt=paths=source_relative
