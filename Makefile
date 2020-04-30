generate:
	protoc -I proto/ -I${GOPATH}/src --go_out=plugins=grpc:proto proto/page.proto
	protoc -I proto/ -I${GOPATH}/src --grpc-gateway_out=logtostderr=true:proto proto/page.proto
	protoc -I proto/ -I${GOPATH}/src --swagger_out=logtostderr=true:proto proto/page.proto
	protoc -I proto/ -I${GOPATH}/src --govalidators_out=logtostderr=true:proto proto/page.proto


install:
	go build . && ./PageService
