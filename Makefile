all: generate

fmt:
	go fmt ./...

install-deps:
	go get -u github.com/jteeuwen/go-bindata/...
	go get -u github.com/elazarl/go-bindata-assetfs/...
	go get -u google.golang.org/grpc
	go get -u github.com/gogo/protobuf/proto
	go get -u github.com/gogo/protobuf/gogoproto
	go get -u github.com/golang/protobuf/protoc-gen-go
	go get -u github.com/gogo/protobuf/protoc-gen-gofast
	go get -u github.com/gogo/protobuf/protoc-gen-gogofaster
	go get -u github.com/gogo/protobuf/protoc-gen-gogoslick
	go get -u -d github.com/grpc-ecosystem/grpc-gateway/protoc-gen-grpc-gateway
	go get -u -d github.com/grpc-ecosystem/grpc-gateway/protoc-gen-swagger
	# git --git-dir=$(GOPATH)/src/github.com/grpc-ecosystem/grpc-gateway/.git --work-tree=$(GOPATH)/src/github.com/grpc-ecosystem/grpc-gateway/ checkout v1.2.2
	go install github.com/grpc-ecosystem/grpc-gateway/protoc-gen-grpc-gateway
	go install github.com/grpc-ecosystem/grpc-gateway/protoc-gen-swagger
	go get -u github.com/go-swagger/go-swagger/cmd/swagger

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
	gofmt -s -w *pb.go *pb.gw.go *pb_test.go swagger.go

generate: generate-proto generate-swagger

generate-swagger: clean-httpapi
	mkdir -p httpapi
	swagger generate server -f dlframework.swagger.json -t httpapi -A dlframework
	swagger generate client -f dlframework.swagger.json -t httpapi -A dlframework
	swagger generate support -f dlframework.swagger.json -t httpapi -A dlframework
	gofmt -s -w httpapi

clean: clean-httpapi
	rm -fr *pb.go *pb.gw.go *pb_test.go swagger.go

clean-httpapi:
	rm -fr httpapi

install-proto:
	./scripts/install-protobuf.sh

travis: install-proto install-deps glide-install logrus-fix generate
	echo "building..."
	go build
