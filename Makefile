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
	mv dlframework.swagger.json dlframework.swagger.json.tmp
	jq -s '.[0] * .[1]' dlframework.swagger.json.tmp swagger-info.json > dlframework.swagger.json
	rm -fr dlframework.swagger.json.tmp
	swagger generate server -f dlframework.swagger.json -t web -A dlframework
	swagger generate client -f dlframework.swagger.json -t web -A dlframework
	swagger generate support -f dlframework.swagger.json -t web -A dlframework

generate-mxnet:
	protoc --plugin=protoc-gen-go=${GOPATH}/bin/protoc-gen-go -Iframeworks/mxnet -I$(GOPATH)/src --gogofaster_out=plugins=grpc:frameworks/mxnet frameworks/mxnet/mxnet.proto
	go-bindata -nomemcopy -prefix frameworks/mxnet/builtin_models/ -pkg mxnet -o frameworks/mxnet/builtin_models_static.go -ignore=.DS_Store -ignore=README.md frameworks/mxnet/builtin_models/...

generate-tensorflow:
	rm -fr frameworks/tensorflow/*pb.go
	protoc --plugin=protoc-gen-go=${GOPATH}/bin/protoc-gen-go -Iframeworks/tensorflow/proto --gogofaster_out=plugins=grpc:frameworks/tensorflow frameworks/tensorflow/proto/allocation_description.proto frameworks/tensorflow/proto/attr_value.proto frameworks/tensorflow/proto/cost_graph.proto frameworks/tensorflow/proto/device_attributes.proto frameworks/tensorflow/proto/function.proto frameworks/tensorflow/proto/graph.proto frameworks/tensorflow/proto/graph_transfer_info.proto frameworks/tensorflow/proto/kernel_def.proto frameworks/tensorflow/proto/log_memory.proto frameworks/tensorflow/proto/node_def.proto frameworks/tensorflow/proto/op_def.proto frameworks/tensorflow/proto/op_gen_overrides.proto frameworks/tensorflow/proto/reader_base.proto frameworks/tensorflow/proto/remote_fused_graph_execute_info.proto frameworks/tensorflow/proto/resource_handle.proto frameworks/tensorflow/proto/step_stats.proto frameworks/tensorflow/proto/summary.proto frameworks/tensorflow/proto/tensor_description.proto frameworks/tensorflow/proto/tensor.proto frameworks/tensorflow/proto/tensor_shape.proto frameworks/tensorflow/proto/tensor_slice.proto frameworks/tensorflow/proto/types.proto frameworks/tensorflow/proto/variable.proto frameworks/tensorflow/proto/versions.proto
	go-bindata -nomemcopy -prefix frameworks/tensorflow/builtin_models/ -pkg tensorflow -o frameworks/tensorflow/builtin_models_static.go -ignore=.DS_Store  -ignore=README.md frameworks/tensorflow/builtin_models/...

linux-brew:
	test -d $HOME/.linuxbrew/bin || git clone https://github.com/Linuxbrew/brew.git $HOME/.linuxbrew
	PATH="$HOME/.linuxbrew/bin:$PATH"
	echo 'export PATH="$HOME/.linuxbrew/bin:$PATH"' >>~/.bash_profile
	export MANPATH="$(brew --prefix)/share/man:$MANPATH"
	export INFOPATH="$(brew --prefix)/share/info:$INFOPATH"
	brew --version
  # Install Buck
	brew tap facebook/fb
	brew install buck
	buck --version

install-proto:
	./scripts/install-protobuf.sh

travis: install-proto install-deps glide-install logrus-fix generate
	echo "building..."
	go build
