#!/bin/bash 
protoc ./cs.proto --go_out=plugins=grpc:./
#mv ./cs.pb.go ../../src/common/cs/
