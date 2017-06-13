# DLFramework

## Model Manifest Format

~~~
name: InceptionNet # name of your model
framework: MXNets
kind: CNN
dataset: ImageNet
container:
  amd64:
    gpu: raiproject/carml-mxnet:amd64-cpu
    cpu: raiproject/carml-mxnet:amd64-gpu
  ppc64le:
    cpu: raiproject/carml-mxnet:ppc64le-gpu
    gpu: raiproject/carml-mxnet:ppc64le-gpu
description: >
  An image-classification convolutional network.
  Inception achieves 21.2% top-1 and 5.6% top-5 error on the ILSVRC 2012 validation dataset.
  It consists of fewer than 25M parameters.
references: 
  - https://arxiv.org/pdf/1512.00567.pdf
licence: MIT
inputs:
  - type: image
    description: the input image
    parameters:
      dimensions: [1, 3, 224, 224]
output:
  type: feature
  parameters:
    features_url: http://data.dmlc.ml/mxnet/models/imagenet/synset.txt
before_preprocess: >
  code... 
preprocess: >
  code... 
after_preprocess: >
  code... 
before_postprocess: >
  code... 
postprocess: >
  code... 
after_postprocess: >
  code... 
model:
  base_url: http://data.dmlc.ml/models/imagenet/inception-bn/
  graph_path: Inception-BN-symbol.json
  weights_path: Inception-BN-0126.params
  is_archive: false 

~~~
