#!/bin/bash

# Checks if a .proto file was provided as an argument
if [ $# -ne 1 ]; then
    echo "Usage: $0 <proto_file>"
    echo "Example: $0 user.proto"
    exit 1
fi

# get a input .proto file name
PROTO_FILE=$1

# Check if a file exists
if [ ! -f "$PROTO_FILE" ]; then
    echo "Error: File $PROTO_FILE not exist"
    exit 1
fi

# execute protoc command
protoc --go_out=. --go_opt=paths=source_relative --go-grpc_out=. --go-grpc_opt=paths=source_relative "$PROTO_FILE"

# Or copy execute protoc command
# protoc --go_out=. --go_opt=paths=source_relative --go-grpc_out=. --go-grpc_opt=paths=source_relative abnormal.proto


# Check whether the command was successfully executed
if [ $? -eq 0 ]; then
    echo "Successfully generated code: $PROTO_FILE"
else
    echo "Error: Code generation failed"
    exit 1
fi