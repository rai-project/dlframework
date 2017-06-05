buildlocalppc:
	docker build . -f Dockerfile.ppc64le -t raiproject/carml-mxnet-ppc64le
buildlocalx64:
	docker build . -f Dockerfile.amd64 -t raiproject/carml-mxnet-amd64

buildlocal: buildlocalppc buildlocalx64
