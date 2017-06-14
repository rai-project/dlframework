all: generate

fmt:
	go fmt ./...

install-deps:
	go get github.com/jteeuwen/go-bindata/...
	go get github.com/elazarl/go-bindata-assetfs/...
	go get github.com/golang/protobuf/{proto,protoc-gen-go}
	go get google.golang.org/grpc
	go get github.com/gogo/protobuf/{proto,gogoproto,protoc-gen-gofast,protoc-gen-gogofaster,protoc-gen-gogoslick}
	go get github.com/grpc-ecosystem/grpc-gateway/{protoc-gen-swagger,protoc-gen-grpc-gateway}
	go get -u github.com/grpc-ecosystem/grpc-gateway/protoc-gen-grpc-gateway
	go get -u github.com/grpc-ecosystem/grpc-gateway/protoc-gen-swagger
	go get -u github.com/golang/protobuf/protoc-gen-go

glide-install:
	glide install --force

logrus-fix:
	rm -fr vendor/github.com/Sirupsen
	find vendor -type f -exec sed -i 's/Sirupsen/sirupsen/g' {} +

generate:
	protoc --plugin=protoc-gen-go=${GOPATH}/bin/protoc-gen-go --proto_path=../../..:. -I$(GOPATH)/src -I$(GOPATH)/src/github.com/grpc-ecosystem/grpc-gateway/third_party/googleapis --swagger_out=logtostderr=true:. --gogofaster_out=plugins=grpc:. dlframework.proto
	jq -s '.[0] * .[1]' dlframework.swagger.json swagger-info.json > dlframework.versioned.swagger.json
	swagger generate server -f dlframework.versioned.swagger.json -t web -A dlframework
	swagger generate client -f dlframework.versioned.swagger.json -t web -A dlframework
	swagger generate support -f dlframework.versioned.swagger.json -t web -A dlframework
