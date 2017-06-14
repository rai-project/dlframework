//go:generate protoc --plugin=protoc-gen-go=${GOPATH}/bin/protoc-gen-go --proto_path=../../..:. -I$GOPATH/src -I$GOPATH/src/github.com/grpc-ecosystem/grpc-gateway/third_party/googleapis --swagger_out=logtostderr=true:. --gogofaster_out=plugins=grpc:. dlframework.proto
//go:generate swagger validate dlframework.swagger.json
//go:generate swagger generate server -f dlframework.swagger.json -t web -A dlframework
//go:generate swagger generate client -f dlframework.swagger.json -t web -A dlframework
//go:generate swagger generate support -f dlframework.swagger.json -t web -A dlframework
package dlframework
