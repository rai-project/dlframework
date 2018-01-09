workspace(name = "dlframework")

git_repository(
    name = "io_bazel_rules_go",
    commit = "570488593c55ad61a18c3d6095344f25da8a84e1",  # master on 2018-01-04
    remote = "https://github.com/bazelbuild/rules_go",
)

git_repository(
    name = "bazel_gazelle",
    commit = "9e43c85089c3247fece397f95dabc1cb63096a59",  # master on 2018-01-09
    remote = "https://github.com/bazelbuild/bazel_gazelle",
)

load("@io_bazel_rules_go//go:def.bzl", "go_rules_dependencies", "go_register_toolchains", "go_repository")
load("@io_bazel_rules_go//proto:def.bzl", "proto_register_toolchains")
load("@bazel_gazelle//:deps.bzl", "gazelle_dependencies")

go_repository(
    name = "com_github_jteeuwen_go_bindata",
    commit = "a0ff2567cfb70903282db057e799fd826784d41d",
    importpath = "github.com/jteeuwen/go-bindata",
)

go_repository(
    name = "com_github_elazarl_go_bindata_assetfs",
    commit = "30f82fa23fd844bd5bb1e5f216db87fd77b5eb43",
    importpath = "github.com/elazarl/go-bindata-assetfs",
)

go_repository(
    name = "org_golang_google_grpc",
    commit = "6913ad5caedced5d627918609375b057963334a5",
    importpath = "google.golang.org/grpc",
)

go_repository(
    name = "com_github_gogo_protobuf_proto",
    commit = "160de10b2537169b5ae3e7e221d28269ef40d311",
    importpath = "github.com/gogo/protobuf/proto",
)

go_repository(
    name = "com_github_gogo_protobuf_gogoproto",
    commit = "160de10b2537169b5ae3e7e221d28269ef40d311",
    importpath = "github.com/gogo/protobuf/gogoproto",
)

go_repository(
    name = "com_github_golang_protobuf_protoc_gen_go",
    commit = "1e59b77b52bf8e4b449a57e6f79f21226d571845",
    importpath = "github.com/golang/protobuf/protoc-gen-go",
)

go_repository(
    name = "com_github_gogo_protobuf_protoc_gen_gofast",
    commit = "160de10b2537169b5ae3e7e221d28269ef40d311",
    importpath = "github.com/gogo/protobuf/protoc-gen-gofast",
)

go_repository(
    name = "com_github_gogo_protobuf_protoc_gen_gogofaster",
    commit = "160de10b2537169b5ae3e7e221d28269ef40d311",
    importpath = "github.com/gogo/protobuf/protoc-gen-gogofaster",
)

go_repository(
    name = "com_github_gogo_protobuf_protoc_gen_gogoslick",
    commit = "160de10b2537169b5ae3e7e221d28269ef40d311",
    importpath = "github.com/gogo/protobuf/protoc-gen-gogoslick",
)

go_repository(
    name = "com_github_grpc_ecosystem_grpc_gateway_protoc_gen_grpc_gateway",
    commit = "61c34cc7e0c7a0d85e4237d665e622640279ff3d",
    importpath = "github.com/grpc-ecosystem/grpc-gateway/protoc-gen-grpc-gateway",
)

go_repository(
    name = "com_github_grpc_ecosystem_grpc_gateway_protoc_gen_swagger",
    commit = "61c34cc7e0c7a0d85e4237d665e622640279ff3d",
    importpath = "github.com/grpc-ecosystem/grpc-gateway/protoc-gen-swagger",
)

go_repository(
    name = "com_github_go_swagger_go_swagger_cmd_swagger",
    commit = "acf3c15f3a1fd86f271220a05558717ec1c61d32",
    importpath = "github.com/go-swagger/go-swagger/cmd/swagger",
)

go_rules_dependencies()

go_register_toolchains()

gazelle_dependencies()
