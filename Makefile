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

generate:
	protoc --plugin=protoc-gen-go=${GOPATH}/bin/protoc-gen-go -I. -I$(GOPATH)/src -I$(GOPATH)/src/github.com/grpc-ecosystem/grpc-gateway/third_party/googleapis --swagger_out=logtostderr=true:. --gogofaster_out=plugins=grpc:. dlframework.proto
	jq -s '.[0] * .[1]' dlframework.swagger.json swagger-info.json > dlframework.versioned.swagger.json
	swagger generate server -f dlframework.versioned.swagger.json -t web -A dlframework
	swagger generate client -f dlframework.versioned.swagger.json -t web -A dlframework
	swagger generate support -f dlframework.versioned.swagger.json -t web -A dlframework

generate-mxnet:
	protoc --plugin=protoc-gen-go=${GOPATH}/bin/protoc-gen-go -Iframeworks/mxnet -I$(GOPATH)/src --gogofaster_out=plugins=grpc:frameworks/mxnet frameworks/mxnet/mxnet.proto
	go-bindata -nomemcopy -prefix frameworks/mxnet/builtin_models/ -pkg mxnet -o frameworks/mxnet/builtin_models_static.go -ignore=.DS_Store frameworks/mxnet/builtin_models/...

install-proto:
	./scripts/install-protobuf.sh

travis: install-proto install-deps glide-install logrus-fix generate
	echo "building..."
	go build
