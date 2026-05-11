# 生成pb命令
protoc --go_out=. --go_opt=paths=source_relative --go-grpc_out=. --go-grpc_opt=paths=source_relative payin.proto

# 只生成数据结构命令
protoc --go_out=. --go_opt=paths=source_relative --go-grpc_ls -laout=. --go-grpc_opt=paths=source_relative payin.proto