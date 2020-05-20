generate:
	protoc -I proto/ -I${GOPATH}/src --go_out=plugins=grpc:proto proto/page.proto
	protoc -I proto/ -I${GOPATH}/src --grpc-gateway_out=logtostderr=true:proto proto/page.proto
	protoc -I proto/ -I${GOPATH}/src --swagger_out=logtostderr=true:proto proto/page.proto
	protoc -I proto/ -I${GOPATH}/src --govalidators_out=logtostderr=true:proto proto/page.proto



install:
	go build . && ./PageService


#go 2 proto dir wer all protos r der n thn run dis cmd

#protoc -I. -I${GOPATH}/src/ -I${GOPATH}/src/github.com/grpc-ecosystem/grpc-gateway/third_party/googleapis -I${GOPATH}/src/github.com/amsokol/protoc-gen-gotag  --gotag_out=xxx="bson+\"-\"":. page.proto
