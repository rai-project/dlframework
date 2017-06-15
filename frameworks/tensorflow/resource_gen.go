//go:generate protoc --plugin=protoc-gen-go=${GOPATH}/bin/protoc-gen-go --proto_path=./proto --gogofaster_out=plugins=grpc:. proto/allocation_description.proto proto/attr_value.proto proto/cost_graph.proto proto/device_attributes.proto proto/function.proto proto/graph.proto proto/graph_transfer_info.proto proto/kernel_def.proto proto/log_memory.proto proto/node_def.proto proto/op_def.proto proto/op_gen_overrides.proto proto/reader_base.proto proto/remote_fused_graph_execute_info.proto proto/resource_handle.proto proto/step_stats.proto proto/summary.proto proto/tensor_description.proto proto/tensor.proto proto/tensor_shape.proto proto/tensor_slice.proto proto/types.proto proto/variable.proto proto/versions.proto
//go:generate go-bindata -nomemcopy -prefix builtin_models/ -pkg tensorflow -o builtin_models_static.go -ignore=.DS_Store ./builtin_models/...
package tensorflow