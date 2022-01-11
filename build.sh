#!/bin/bash
######################################################################
#@编译脚本
#@Date 2022-01-10 
#@Author Tom.Chen 
######################################################################
function buildProto()
{
	protoc ./pb_proto/common/cs/cs.proto --go_out=plugins=grpc:./	
	mv *.pb.go ./cs
}
function main()
{
	buildProto
	make -j4 
	rm -f build
	mv mgr_service ./bin
}
main 
