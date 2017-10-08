#!/bin/bash

FRAMEWORK=MxNet
FRAMEWORK_VERSION=0.11.0
BATCH=1024
URLFILE=urlsfile
OUTDIR=mxnet_out

declare -a array1=( BVLC-AlexNet BVLC-GoogLeNet BVLC-Reference-CaffeNet \
    BVLC-Reference-RCNN-ILSVRC13 InceptionBN-21K \
    ResNeXt101-32x4d ResNeXt26-32x4d-priv ResNeXt50-32x4d \
    ResNet101 ResNet152 ResNet50 SqueezeNet \
    VGG16 VGG16_SOD VGG19 WRN50-2 )

declare -a array2=( ResNet101 ) # v2.0
declare -a array3=( Inception) # v3.0
declare -a array4=( Inception ) # v4.0
declare -a array5=( SqueezeNet ) # v1.1

rm -rf $OUTDIR

for i in "${array1[@]}"
do
    echo $i
    ./run --modelName $i --modelVersion "1.0" collect $URLFILE $BATCH $OUTDIR
done

for i in "${array2[@]}"
do
    echo $i
    ./run --modelName $i --modelVersion "2.0" collect $URLFILE $BATCH $OUTDIR
done

for i in "${array3[@]}"
do
    echo $i
    ./run --modelName $i --modelVersion "3.0" collect $URLFILE $BATCH $OUTDIR
done

for i in "${array4[@]}"
do
    echo $i
    ./run --modelName $i --modelVersion "4.0" collect $URLFILE $BATCH $OUTDIR
done

for i in "${array5[@]}"
do
    echo $i
    ./run --modelName $i --modelVersion "1.1" collect $URLFILE $BATCH $OUTDIR
done
