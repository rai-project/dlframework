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
	go get -d github.com/grpc-ecosystem/grpc-gateway/protoc-gen-grpc-gateway
	go get -d github.com/grpc-ecosystem/grpc-gateway/protoc-gen-swagger
	git --git-dir=$(GOPATH)/src/github.com/grpc-ecosystem/grpc-gateway/.git --work-tree=$(GOPATH)/src/github.com/grpc-ecosystem/grpc-gateway/ checkout v1.2.2
	go install github.com/grpc-ecosystem/grpc-gateway/protoc-gen-grpc-gateway
	go install github.com/grpc-ecosystem/grpc-gateway/protoc-gen-swagger
	go get github.com/go-swagger/go-swagger/cmd/swagger

glide-install:
	glide install --force

logrus-fix:
	rm -fr vendor/github.com/Sirupsen
	find vendor -type f -exec sed -i 's/Sirupsen/sirupsen/g' {} +

generate-proto:
	rm -fr swagger.go
	protoc --plugin=protoc-gen-go=${GOPATH}/bin/protoc-gen-go -I. -I$(GOPATH)/src -I$(GOPATH)/src/github.com/golang/protobuf/proto -I$(GOPATH)/src/github.com/grpc-ecosystem/grpc-gateway/third_party/googleapis --gogofaster_out=Mgoogle/protobuf/struct.proto=github.com/gogo/protobuf/types,plugins=grpc:. dlframework.proto
	protoc -I. -I$(GOPATH)/src -I$(GOPATH)/src/github.com/golang/protobuf/proto -I$(GOPATH)/src/github.com/grpc-ecosystem/grpc-gateway/third_party/googleapis --grpc-gateway_out=logtostderr=true:. dlframework.proto
	protoc -I. -I$(GOPATH)/src -I$(GOPATH)/src/github.com/golang/protobuf/proto -I$(GOPATH)/src/github.com/grpc-ecosystem/grpc-gateway/third_party/googleapis --swagger_out=logtostderr=true:. dlframework.proto
	mv dlframework.swagger.json dlframework.swagger.json.tmp
	jq -s '.[0] * .[1]' dlframework.swagger.json.tmp swagger_info.json > dlframework.swagger.json
	rm -fr dlframework.swagger.json.tmp
	go run scripts/includetext.go

generate: generate-proto generate-frameworks generate-swagger

generate-swagger: clean-web
	swagger generate server -f dlframework.swagger.json -t web -A dlframework
	swagger generate client -f dlframework.swagger.json -t web -A dlframework
	swagger generate support -f dlframework.swagger.json -t web -A dlframework

clean: clean-web clean-frameworks
	rm -fr *pb.go *pb.gw.go *pb_test.go swagger.go

clean-web:
	rm -fr web

clean-mxnet:
	rm -fr frameworks/mxnet/builtin_models_static.go frameworks/mxnet/*pb.go

generate-mxnet: clean-mxnet
	protoc --plugin=protoc-gen-go=${GOPATH}/bin/protoc-gen-go -Iframeworks/mxnet -I$(GOPATH)/src --gogofaster_out=plugins=grpc:frameworks/mxnet frameworks/mxnet/mxnet.proto
	go-bindata -nomemcopy -prefix frameworks/mxnet/builtin_models/ -pkg mxnet -o frameworks/mxnet/builtin_models_static.go -ignore=.DS_Store -ignore=README.md frameworks/mxnet/builtin_models/...

clean-tensorflow:
	rm -fr frameworks/tensorflow/builtin_models_static.go frameworks/tensorflow/*pb.go

generate-tensorflow: clean-tensorflow
	protoc --plugin=protoc-gen-go=${GOPATH}/bin/protoc-gen-go -Iframeworks/tensorflow/proto --gogofaster_out=Mgoogle/protobuf/any.proto=github.com/gogo/protobuf/types,plugins=grpc:frameworks/tensorflow frameworks/tensorflow/proto/allocation_description.proto frameworks/tensorflow/proto/attr_value.proto frameworks/tensorflow/proto/cluster.proto frameworks/tensorflow/proto/config.proto frameworks/tensorflow/proto/control_flow.proto frameworks/tensorflow/proto/cost_graph.proto frameworks/tensorflow/proto/debug.proto frameworks/tensorflow/proto/device_attributes.proto frameworks/tensorflow/proto/device_properties.proto frameworks/tensorflow/proto/error_codes.proto frameworks/tensorflow/proto/function.proto frameworks/tensorflow/proto/graph.proto frameworks/tensorflow/proto/graph_transfer_info.proto frameworks/tensorflow/proto/kernel_def.proto frameworks/tensorflow/proto/log_memory.proto frameworks/tensorflow/proto/master.proto frameworks/tensorflow/proto/master_service.proto frameworks/tensorflow/proto/meta_graph.proto frameworks/tensorflow/proto/named_tensor.proto frameworks/tensorflow/proto/node_def.proto frameworks/tensorflow/proto/op_def.proto frameworks/tensorflow/proto/op_gen_overrides.proto frameworks/tensorflow/proto/queue_runner.proto frameworks/tensorflow/proto/reader_base.proto frameworks/tensorflow/proto/remote_fused_graph_execute_info.proto frameworks/tensorflow/proto/resource_handle.proto frameworks/tensorflow/proto/rewriter_config.proto frameworks/tensorflow/proto/saved_model.proto frameworks/tensorflow/proto/saver.proto frameworks/tensorflow/proto/step_stats.proto frameworks/tensorflow/proto/summary.proto frameworks/tensorflow/proto/tensor_bundle.proto frameworks/tensorflow/proto/tensor_description.proto frameworks/tensorflow/proto/tensorflow_server.proto frameworks/tensorflow/proto/tensor.proto frameworks/tensorflow/proto/tensor_shape.proto frameworks/tensorflow/proto/tensor_slice.proto frameworks/tensorflow/proto/types.proto frameworks/tensorflow/proto/variable.proto frameworks/tensorflow/proto/versions.proto frameworks/tensorflow/proto/worker.proto frameworks/tensorflow/proto/worker_service.proto
	go-bindata -nomemcopy -prefix frameworks/tensorflow/builtin_models/ -pkg tensorflow -o frameworks/tensorflow/builtin_models_static.go -ignore=.DS_Store  -ignore=README.md frameworks/tensorflow/builtin_models/...

clean-caffe:
	rm -fr frameworks/caffe/builtin_models_static.go frameworks/caffe/*pb.go

generate-caffe: clean-caffe
	protoc --plugin=protoc-gen-go=${GOPATH}/bin/protoc-gen-go -Iframeworks/caffe/proto --gogofaster_out=plugins=grpc:frameworks/caffe frameworks/caffe/proto/caffe.proto
	go-bindata -nomemcopy -prefix frameworks/caffe/builtin_models/ -pkg tensorflow -o frameworks/caffe/builtin_models_static.go -ignore=.DS_Store  -ignore=README.md frameworks/caffe/builtin_models/...


clean-frameworks: clean-mxnet clean-caffe clean-tensorflow

generate-frameworks: generate-mxnet generate-caffe generate-tensorflow

install-proto:
	./scripts/install-protobuf.sh

travis: install-proto install-deps glide-install logrus-fix generate
	echo "building..."
	go build
