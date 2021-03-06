#!/bin/bash

FRAMEWORK=MxNet
FRAMEWORK_VERSION=0.11.0
BATCH=512
URLFILE=urlsfile
OUTDIR=mxnet_out

declare -a array1=( BVLC-AlexNet BVLC-GoogLeNet BVLC-Reference-CaffeNet \
    BVLC-Reference-RCNN-ILSVRC13 \
    ResNeXt101-32x4d ResNeXt26-32x4d ResNeXt50-32x4d \
    ResNet101 ResNet152 ResNet50 SqueezeNet )

declare -a array2=( ResNet101 InceptionBN-21K WRN50 ) # v2.0
declare -a array3=( Inception) # v3.0
declare -a array4=( Inception ) # v4.0
declare -a array5=( SqueezeNet ) # v1.1

# VGG16, VGG16_SOD, VGG19 crashe for batch = 2, 4
declare -a array6=( VGG16 VGG16_SOD VGG19  )

for i in "${array1[@]}"
do
    echo $i
    ./run --modelName $i --modelVersion "1.0" collect $URLFILE $BATCH $OUTDIR/$i-v1.0
done

# for i in "${array2[@]}"
# do
#     echo $i
#     ./run --modelName $i --modelVersion "2.0" collect $URLFILE $BATCH $OUTDIR/$i-v2.0
# done

# for i in "${array3[@]}"
# do
#     echo $i
#     ./run --modelName $i --modelVersion "3.0" collect $URLFILE $BATCH $OUTDIR/$i-v3.0
# done

# for i in "${array4[@]}"
# do
#     echo $i
#     ./run --modelName $i --modelVersion "4.0" collect $URLFILE $BATCH $OUTDIR/$i-v4.0
# done

# for i in "${array5[@]}"
# do
#     echo $i
#     ./run --modelName $i --modelVersion "1.1" collect $URLFILE $BATCH $OUTDIR/$i-v1.1
# done
