protoc --go_out=. --go_opt=paths=source_relative --go-grpc_out=. --go-grpc_opt=paths=source_relative settlement.proto

# 只生成数据结构命令
protoc --go_out=. --go_opt=paths=source_relative --go-grpc_out=. --go-grpc_opt=paths=source_relative xxx.proto

protoc --go_out=. --go_opt=paths=source_relative --go-grpc_out=. --go-grpc_opt=paths=source_relative --validate_out="lang=go:." xxx.proto

# ============
protoc --proto_path=. --proto_path=./rpc_api --go_out=. --go_opt=paths=source_relative --go-grpc_out=. --go-grpc_opt=paths=source_relative --validate_out="lang=go:." example.proto

protoc --proto_path=. --proto_path=./rpc_api --go_out=. --validate_out="lang=go:." example.proto