all: buildlocalbvlc buildlocalibm buildlocalnv

buildlocalbvlc:
	docker build -f Dockerfile.bvlc -t raiproject/powerai:caffe-bvlc .

buildlocalibm:
	docker build -f Dockerfile.ibm -t raiproject/powerai:caffe-ibm .

buildlocalnv:
	docker build -f Dockerfile.nv -t raiproject/powerai:caffe-nv .

