all: generate

fmt:
	go fmt ./...

install-deps:
	go get github.com/jteeuwen/go-bindata/...
	go get github.com/elazarl/go-bindata-assetfs/...
	go get google.golang.org/grpc
	go get github.com/gogo/protobuf/proto
	go get github.com/gogo/protobuf/gogoproto
	go get github.com/golang/protobuf/protoc-gen-go
	go get github.com/gogo/protobuf/protoc-gen-gofast
	go get github.com/gogo/protobuf/protoc-gen-gogofaster
	go get github.com/gogo/protobuf/protoc-gen-gogoslick
	go get github.com/grpc-ecosystem/grpc-gateway/protoc-gen-swagger
	go get github.com/grpc-ecosystem/grpc-gateway/protoc-gen-grpc-gateway
	go get github.com/grpc-ecosystem/grpc-gateway/protoc-gen-swagger
	go get github.com/go-swagger/go-swagger/cmd/swagger

glide-install:
	glide install --force

logrus-fix:
	rm -fr vendor/github.com/Sirupsen
	find vendor -type f -exec sed -i 's/Sirupsen/sirupsen/g' {} +

generate-proto:
	protoc --plugin=protoc-gen-go=${GOPATH}/bin/protoc-gen-go -I. -I$(GOPATH)/src -I$(GOPATH)/src/github.com/golang/protobuf/proto -I$(GOPATH)/src/github.com/grpc-ecosystem/grpc-gateway/third_party/googleapis --swagger_out=logtostderr=true:. --gogofaster_out=Mgoogle/protobuf/struct.proto=github.com/gogo/protobuf/types,plugins=grpc:. dlframework.proto

generate: generate-proto
	jq -s '.[0] * .[1]' dlframework.swagger.json swagger-info.json > dlframework.versioned.swagger.json
	swagger generate server -f dlframework.versioned.swagger.json -t web -A dlframework
	swagger generate client -f dlframework.versioned.swagger.json -t web -A dlframework
	swagger generate support -f dlframework.versioned.swagger.json -t web -A dlframework

generate-mxnet:
	protoc --plugin=protoc-gen-go=${GOPATH}/bin/protoc-gen-go -Iframeworks/mxnet -I$(GOPATH)/src --gogofaster_out=plugins=grpc:frameworks/mxnet frameworks/mxnet/mxnet.proto
	go-bindata -nomemcopy -prefix frameworks/mxnet/builtin_models/ -pkg mxnet -o frameworks/mxnet/builtin_models_static.go -ignore=.DS_Store frameworks/mxnet/builtin_models/...

generate-tensorflow:
	protoc --plugin=protoc-gen-go=${GOPATH}/bin/protoc-gen-go -Iframeworks/tensorflow --proto_path=./proto --gogofaster_out=plugins=grpc:frameworks/tensorflow proto/allocation_description.proto proto/attr_value.proto proto/cost_graph.proto proto/device_attributes.proto proto/function.proto proto/graph.proto proto/graph_transfer_info.proto proto/kernel_def.proto proto/log_memory.proto proto/node_def.proto proto/op_def.proto proto/op_gen_overrides.proto proto/reader_base.proto proto/remote_fused_graph_execute_info.proto proto/resource_handle.proto proto/step_stats.proto proto/summary.proto proto/tensor_description.proto proto/tensor.proto proto/tensor_shape.proto proto/tensor_slice.proto proto/types.proto proto/variable.proto proto/versions.proto
	go-bindata -nomemcopy -prefix frameworks/tensorflow/builtin_models/ -pkg mxnet -o frameworks/tensorflow/builtin_models_static.go -ignore=.DS_Store frameworks/tensorflow/builtin_models/...


install-proto:
	./scripts/install-protobuf.sh

travis: install-proto install-deps glide-install logrus-fix generate
	echo "building..."
	go build
