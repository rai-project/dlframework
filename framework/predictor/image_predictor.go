package predictor

import (
	"bufio"
	"bytes"
	"context"
	"image"
	"image/color"
	"image/jpeg"
	"os"
	"path/filepath"
	"sort"
	"strconv"

	"github.com/k0kubun/pp"
	"github.com/pkg/errors"
	"github.com/rai-project/dlframework"
	"github.com/rai-project/dlframework/framework/feature"
	raiimage "github.com/rai-project/image"
	"github.com/rai-project/image/types"
	imageTypes "github.com/rai-project/image/types"
	"github.com/rai-project/utils"
	"github.com/spf13/cast"
	yaml "gopkg.in/yaml.v2"
	"gorgonia.org/tensor"
)

type Method int

const (
	TopLeft Method = iota
	Center
	InvalidCropMethod Method = 9999
)

type PreprocessOptions struct {
	ElementType     string
	MeanImage       []float32
	Dims            []int
	MaxDimension    *int
	KeepAspectRatio *bool
	Scale           float32
	ColorMode       types.Mode
	Layout          raiimage.Layout
	CropMethod      Method
	CropRatio       float32
}

type ImagePredictor struct {
	Base
	Metadata map[string]interface{}
}

func (p ImagePredictor) GetInputLayerName(layer string) (string, error) {
	model := p.Model
	modelInputs := model.GetInputs()
	for _, input := range modelInputs {
		typeParameters := input.GetParameters()
		name, err := p.GetTypeParameter(typeParameters, layer)
		if err != nil {
			return "", err
		}
		return name, nil
	}
	return "", errors.New("cannot find input layer name")
}

func (p ImagePredictor) GetOutputLayerIndex(layer string) (int, error) {
	model := p.Model
	modelOutput := model.GetOutput()
	typeParameters := modelOutput.GetParameters()
	str, err := p.GetTypeParameter(typeParameters, layer)
	if err != nil {
		return 0, err
	}
	index, err := strconv.Atoi(str)
	if err != nil {
		return 0, err
	}
	return index, nil
}

func (p ImagePredictor) GetInputDimensions() ([]int, error) {
	model := p.Model
	modelInputs := model.GetInputs()

	typeParameters := modelInputs[0].GetParameters()

	if typeParameters == nil {
		return nil, errors.New("invalid type parameters")
	}
	pdims, ok := typeParameters["dimensions"]
	if !ok {
		log.Debug("arbitrary input dimensions")
		return nil, nil
	}
	pdimsVal := pdims.GetValue()
	if pdimsVal == "" {
		return nil, errors.New("invalid input dimensions")
	}

	var dims []int
	if err := yaml.Unmarshal([]byte(pdimsVal), &dims); err != nil {
		return nil, errors.Errorf("unable to get input dimensions %v as an integer slice", pdimsVal)
	}
	if len(dims) == 1 {
		dims = []int{3, dims[0], dims[0]}
	}
	if len(dims) > 3 {
		return nil, errors.Errorf("expecting a dimensions size of 1 or 3, but got %v. do not put the batch size in the input dimensions.", len(dims))
	}

	return dims, nil
}

func (p ImagePredictor) GetMeanPath() string {
	model := p.Model
	return dlframework.CleanString(filepath.Join(p.WorkDir, model.GetName()+".mean"))
}

func (p ImagePredictor) GetMeanImage() ([]float32, error) {
	model := p.Model
	modelInputs := model.GetInputs()
	typeParameters := modelInputs[0].GetParameters()
	if typeParameters == nil {
		return nil, errors.New("invalid type parameters")
	}
	pmean, ok := typeParameters["mean"]
	if !ok {
		log.Debug("using 0,0,0 as the mean image")
		return []float32{0, 0, 0}, nil
	}

	pmeanVal := pmean.GetValue()
	if pmeanVal == "" {
		return nil, errors.New("invalid mean image")
	}

	var vals []float32
	if err := yaml.Unmarshal([]byte(pmeanVal), &vals); err == nil {
		return vals, nil
	}
	var val float32
	if err := yaml.Unmarshal([]byte(pmeanVal), &val); err != nil {
		return nil, errors.Errorf("unable to get image mean %v as a float or slice", pmeanVal)
	}

	return []float32{val, val, val}, nil
}

func (p ImagePredictor) GetScale() (float32, error) {
	model := p.Model
	modelInputs := model.GetInputs()
	if len(modelInputs) == 0 {
		return 1.0, nil
	}
	typeParameters := modelInputs[0].GetParameters()
	if typeParameters == nil {
		return 1.0, errors.New("invalid type parameters")
	}
	pscale, ok := typeParameters["scale"]
	if !ok {
		// log.Debug("no scaling")
		return 1.0, nil
	}
	pscaleVal := pscale.GetValue()
	if pscaleVal == "" {
		return 1.0, nil
	}

	var val float32
	if err := yaml.Unmarshal([]byte(pscaleVal), &val); err != nil {
		return 1.0, errors.Errorf("unable to get scale %v as a float", pscaleVal)
	}

	return val, nil
}

func (p ImagePredictor) GetMaxDimension() (int, error) {
	model := p.Model
	modelInputs := model.GetInputs()
	if len(modelInputs) == 0 {
		return 0, errors.New("no inputs")
	}
	typeParameters := modelInputs[0].GetParameters()
	if typeParameters == nil {
		return 0, errors.New("invalid type parameters")
	}
	pscale, ok := typeParameters["max_dimension"]
	if !ok {
		return 0, errors.New("no max dimension")
	}
	pscaleVal := pscale.GetValue()
	if pscaleVal == "" {
		return 0, errors.New("no max dimension value")
	}

	var val int
	if err := yaml.Unmarshal([]byte(pscaleVal), &val); err != nil {
		return 0, errors.Errorf("unable to get max dimension %v as a int", pscaleVal)
	}

	return val, nil
}

func (p ImagePredictor) GetKeepAspectRatio() (bool, error) {
	model := p.Model
	modelInputs := model.GetInputs()
	if len(modelInputs) == 0 {
		return false, errors.New("no inputs")
	}
	typeParameters := modelInputs[0].GetParameters()
	if typeParameters == nil {
		return false, errors.New("invalid type parameters")
	}
	pscale, ok := typeParameters["keep_aspect_ratio"]
	if !ok {
		return false, errors.New("no keep aspect ratio")
	}
	pscaleVal := pscale.GetValue()
	if pscaleVal == "" {
		return false, errors.New("no keep aspect ratio value")
	}

	var val bool
	if err := yaml.Unmarshal([]byte(pscaleVal), &val); err != nil {
		return false, errors.Errorf("unable to get keep aspect ratio %v as a bool", pscaleVal)
	}

	return val, nil
}

func (p ImagePredictor) GetLayout(defaultLayout raiimage.Layout) raiimage.Layout {
	model := p.Model
	modelInputs := model.GetInputs()
	typeParameters := modelInputs[0].GetParameters()
	if typeParameters == nil {
		return defaultLayout
	}
	pscale, ok := typeParameters["layout"]
	if !ok {
		return defaultLayout
	}
	pscaleVal := pscale.GetValue()
	if pscaleVal == "" {
		return defaultLayout
	}

	var val string
	if err := yaml.Unmarshal([]byte(pscaleVal), &val); err != nil {
		log.Errorf("unable to get color_mode %v as a string", pscaleVal)
		return defaultLayout
	}

	switch val {
	case "CHW":
		return raiimage.CHWLayout
	case "HWC":
		return raiimage.HWCLayout
	default:
		log.Error("invalid image mode specified " + val)
		return raiimage.InvalidLayout
	}
}

func (p ImagePredictor) GetColorMode(defaultMode types.Mode) types.Mode {
	model := p.Model
	modelInputs := model.GetInputs()
	typeParameters := modelInputs[0].GetParameters()
	if typeParameters == nil {
		return defaultMode
	}
	pscale, ok := typeParameters["color_mode"]
	if !ok {
		return defaultMode
	}
	pscaleVal := pscale.GetValue()
	if pscaleVal == "" {
		return defaultMode
	}

	var val string
	if err := yaml.Unmarshal([]byte(pscaleVal), &val); err != nil {
		log.Errorf("unable to get color_mode %v as a string", pscaleVal)
		return defaultMode
	}

	switch val {
	case "RGB":
		return types.RGBMode
	case "BGR":
		return types.BGRMode
	default:
		log.Error("invalid image mode specified " + val)
		return types.InvalidMode
	}
}

func (p ImagePredictor) GetCropMethod(defaultMethod Method) Method {
	model := p.Model
	modelInputs := model.GetInputs()
	typeParameters := modelInputs[0].GetParameters()
	if typeParameters == nil {
		return defaultMethod
	}
	pcropMethod, ok := typeParameters["crop_method"]
	if !ok {
		return defaultMethod
	}
	pcropMethodVal := pcropMethod.GetValue()
	if pcropMethodVal == "" {
		return defaultMethod
	}

	var val string
	if err := yaml.Unmarshal([]byte(pcropMethodVal), &val); err != nil {
		log.Errorf("unable to get color_mode %v as a string", pcropMethodVal)
		return defaultMethod
	}

	switch val {
	case "topleft":
		return TopLeft
	case "center":
		return Center
	default:
		log.Error("invalid image mode specified " + val)
		return InvalidCropMethod
	}
}

func (p ImagePredictor) GetCropRatio(defaultCropRatio float32) float32 {
	model := p.Model
	modelInputs := model.GetInputs()
	typeParameters := modelInputs[0].GetParameters()
	if typeParameters == nil {
		return defaultCropRatio
	}
	pCropRatio, ok := typeParameters["crop_ratio"]
	if !ok {
		return defaultCropRatio
	}
	pCropRatioVal := pCropRatio.GetValue()
	if pCropRatioVal == "" {
		return defaultCropRatio
	}

	var val float32
	if err := yaml.Unmarshal([]byte(pCropRatioVal), &val); err != nil {
		log.Errorf("unable to get crop_ratio %v as a float32", pCropRatioVal)
		return defaultCropRatio
	}

	return val
}

func (p ImagePredictor) GetPreprocessOptions() (PreprocessOptions, error) {
	mean, err := p.GetMeanImage()
	if err != nil {
		return PreprocessOptions{}, err
	}
	scale, err := p.GetScale()
	if err != nil {
		return PreprocessOptions{}, err
	}

	imageDims, err := p.GetInputDimensions()
	if err != nil {
		imageDims = nil
	}

	maxDim0, err := p.GetMaxDimension()
	maxDim := &maxDim0
	if err != nil {
		maxDim = nil
	}
	keepAspectRatio0, err := p.GetKeepAspectRatio()
	keepAspectRatio := &keepAspectRatio0
	if err != nil {
		keepAspectRatio = nil
	}

	preprocOpts := PreprocessOptions{
		ElementType:     p.Model.GetElementType(),
		MeanImage:       mean,
		Scale:           scale,
		Dims:            imageDims,
		MaxDimension:    maxDim,
		KeepAspectRatio: keepAspectRatio,
		ColorMode:       p.GetColorMode(imageTypes.RGBMode),
		Layout:          p.GetLayout(raiimage.HWCLayout),
		CropMethod:      p.GetCropMethod(Center),
		CropRatio:       p.GetCropRatio(1.0),
	}

	return preprocOpts, nil
}

func (p ImagePredictor) GetLabels() ([]string, error) {
	var labels []string
	f, err := os.Open(p.GetFeaturesPath())
	if err != nil {
		return nil, errors.Wrapf(err, "cannot read %s", p.GetFeaturesPath())
	}
	defer f.Close()
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := scanner.Text()
		labels = append(labels, line)
	}
	return labels, nil
}

func (p ImagePredictor) iCreateClassificationFeatures2DSlice(ctx context.Context, probabilities [][]float32, labels []string) ([]dlframework.Features, error) {
	batchSize := p.BatchSize()
	if len(probabilities) < 1 {
		return nil, errors.New("len(probabilities) < 1")
	}
	featureLen := len(probabilities[0])
	features := make([]dlframework.Features, batchSize)

	for ii := 0; ii < batchSize; ii++ {
		rprobs := make([]*dlframework.Feature, featureLen)
		for jj := 0; jj < featureLen; jj++ {
			rprobs[jj] = feature.New(
				feature.ClassificationIndex(int32(jj)),
				feature.ClassificationLabel(labels[jj]),
				feature.Probability(probabilities[ii][jj]),
			)
		}
		sort.Sort(dlframework.Features(rprobs))
		features[ii] = rprobs
	}

	return features, nil
}

func (p ImagePredictor) CreateClassificationFeaturesFrom1D(ctx context.Context, probabilities []float32, labels []string) ([]dlframework.Features, error) {

	batchSize := p.BatchSize()
	featureLen := len(probabilities) / batchSize

	probs := tensor.New(
		tensor.Of(tensor.Float32),
		tensor.WithBacking(probabilities),
		tensor.WithShape(batchSize, featureLen),
	)

	return p.CreateClassificationFeatures(ctx, probs, labels)
}

func (p ImagePredictor) CreateClassificationFeatures(ctx context.Context, probabilities0 interface{}, labels []string) ([]dlframework.Features, error) {
	if slc, ok := probabilities0.([][]float32); ok {
		return p.iCreateClassificationFeatures2DSlice(ctx, slc, labels)
	}

	probabilities, ok := probabilities0.(tensor.Tensor)
	if !ok {
		return nil, errors.New("expecting an input tensor")
	}

	batchSize := p.BatchSize()
	if probabilities.Size() == 0 {
		return nil, errors.New("len(probabilities) == 0")
	}
	if probabilities.Dtype() == tensor.Float32 {
		return nil, errors.New("invalid data type")
	}

	featureLen := probabilities.Shape()[0]
	features := make([]dlframework.Features, batchSize)

	for ii := 0; ii < batchSize; ii++ {
		rprobs := make([]*dlframework.Feature, featureLen)
		for jj := 0; jj < featureLen; jj++ {
			iprob, err := probabilities.At(ii, jj)
			if err != nil {
				return nil, errors.Wrapf(err, "invalid probability value at (%v,%v)", ii, jj)
			}
			prob := cast.ToFloat32(iprob)
			rprobs[jj] = feature.New(
				feature.ClassificationIndex(int32(jj)),
				feature.ClassificationLabel(labels[jj]),
				feature.Probability(prob),
			)
		}
		sort.Sort(dlframework.Features(rprobs))
		features[ii] = rprobs
	}

	return features, nil
}

func (p ImagePredictor) iCreateBoundingBoxFeaturesSlice(ctx context.Context, probabilities [][]float32, classes [][]float32, boxes [][][]float32, labels []string) ([]dlframework.Features, error) {
	batchSize := p.BatchSize()
	if len(probabilities) < 1 {
		return nil, errors.New("len(probabilities) < 1")
	}
	featureLen := len(probabilities[0])
	features := make([]dlframework.Features, batchSize)

	for ii := 0; ii < batchSize; ii++ {
		rprobs := make([]*dlframework.Feature, featureLen)
		for jj := 0; jj < featureLen; jj++ {
			rprobs[jj] = feature.New(
				feature.BoundingBoxType(),
				feature.BoundingBoxXmin(boxes[ii][jj][1]),
				feature.BoundingBoxXmax(boxes[ii][jj][3]),
				feature.BoundingBoxYmin(boxes[ii][jj][0]),
				feature.BoundingBoxYmax(boxes[ii][jj][2]),
				feature.BoundingBoxIndex(int32(classes[ii][jj])),
				feature.BoundingBoxLabel(labels[int32(classes[ii][jj])]),
				feature.Probability(probabilities[ii][jj]),
			)
		}
		sort.Sort(dlframework.Features(rprobs))
		features[ii] = rprobs
	}

	return features, nil
}

func (p ImagePredictor) CreateBoundingBoxFeatures(ctx context.Context, probabilities0 interface{}, classes0 interface{}, boxes0 interface{}, labels []string) ([]dlframework.Features, error) {
	if slc, ok := probabilities0.([][]float32); ok {
		return p.iCreateBoundingBoxFeaturesSlice(ctx, slc, classes0.([][]float32), boxes0.([][][]float32), labels)
	}

	probabilities, ok := probabilities0.(tensor.Tensor)
	if !ok {
		return nil, errors.New("expecting an input probabilities tensor")
	}

	classes, ok := classes0.(tensor.Tensor)
	if !ok {
		return nil, errors.New("expecting an input classes tensor")
	}

	boxes, ok := boxes0.(tensor.Tensor)
	if !ok {
		return nil, errors.New("expecting an input boxes tensor")
	}

	batchSize := p.BatchSize()
	if probabilities.Size() == 0 {
		return nil, errors.New("len(probabilities) < 1")
	}

	featureLen := probabilities.Shape()[0]
	features := make([]dlframework.Features, batchSize)

	for ii := 0; ii < batchSize; ii++ {
		rprobs := make([]*dlframework.Feature, featureLen)
		for jj := 0; jj < featureLen; jj++ {
			iclass, err := classes.At(ii, jj)
			if err != nil {
				return nil, errors.Wrapf(err, "invalid class value at (%v,%v)", ii, jj)
			}
			iprob, err := probabilities.At(ii, jj)
			if err != nil {
				return nil, errors.Wrapf(err, "invalid probability value at (%v,%v)", ii, jj)
			}
			iboxYMin, err := boxes.At(ii, jj, 0)
			if err != nil {
				return nil, errors.Wrapf(err, "invalid box y min value at (%v,%v, 0)", ii, jj)
			}
			iboxXMin, err := boxes.At(ii, jj, 1)
			if err != nil {
				return nil, errors.Wrapf(err, "invalid box x min value at (%v,%v, 1)", ii, jj)
			}
			iboxYMax, err := boxes.At(ii, jj, 2)
			if err != nil {
				return nil, errors.Wrapf(err, "invalid box y max value at (%v,%v, 2)", ii, jj)
			}
			iboxXMax, err := boxes.At(ii, jj, 3)
			if err != nil {
				return nil, errors.Wrapf(err, "invalid box x max value at (%v,%v, 3)", ii, jj)
			}
			class := cast.ToInt32(iclass)
			prob := cast.ToFloat32(iprob)
			boxYMin := cast.ToFloat32(iboxYMin)
			boxYMax := cast.ToFloat32(iboxYMax)
			boxXMin := cast.ToFloat32(iboxXMin)
			boxXMax := cast.ToFloat32(iboxXMax)
			rprobs[jj] = feature.New(
				feature.BoundingBoxType(),
				feature.BoundingBoxXmin(boxXMin),
				feature.BoundingBoxXmax(boxXMax),
				feature.BoundingBoxYmin(boxYMin),
				feature.BoundingBoxYmax(boxYMax),
				feature.BoundingBoxIndex(class),
				feature.BoundingBoxLabel(labels[class]),
				feature.Probability(prob),
			)
		}
		sort.Sort(dlframework.Features(rprobs))
		features[ii] = rprobs
	}

	return features, nil
}

func (p ImagePredictor) iCreateSemanticSegmentFeaturesSlice(ctx context.Context, masks [][][]int64, labels []string) ([]dlframework.Features, error) {
	batchSize := p.BatchSize()
	if len(masks) < 1 {
		return nil, errors.New("len(masks) < 1")
	}
	targetHeight := len(masks[0])
	targetWidth := len(masks[0][0])
	features := make([]dlframework.Features, batchSize)
	featureLen := 1
	for ii := 0; ii < batchSize; ii++ {
		rprobs := make([]*dlframework.Feature, featureLen)
		mask := masks[ii]
		for jj := 0; jj < featureLen; jj++ {
			rprobs[jj] = feature.New(
				feature.SemanticSegmentType(),
				feature.SemanticSegmentHeight(int32(targetHeight)),
				feature.SemanticSegmentWidth(int32(targetWidth)),
				feature.SemanticSegmentIntMask(flattenInt32Slice(mask)),
				feature.Probability(1.0),
			)
		}
		features[ii] = rprobs
	}

	return features, nil
}

func (p ImagePredictor) CreateSemanticSegmentFeatures(ctx context.Context, masks0 interface{}, labels []string) ([]dlframework.Features, error) {
	if masks, ok := masks0.([][][]int64); ok {
		return p.iCreateSemanticSegmentFeaturesSlice(ctx, masks, labels)
	}

	masks, ok := masks0.(tensor.Tensor)
	if !ok {
		return nil, errors.New("expecting an input masks tensor")
	}

	batchSize := p.BatchSize()
	if masks.Size() == 0 {
		return nil, errors.New("len(masks) < 1")
	}
	if masks.Dims() == 3 {
		return nil, errors.New("rank(masks) != 3")
	}
	targetHeight := masks.Shape()[1]
	targetWidth := masks.Shape()[2]
	features := make([]dlframework.Features, batchSize)
	featureLen := 1
	for ii := 0; ii < batchSize; ii++ {
		rprobs := make([]*dlframework.Feature, featureLen)
		iMask, err := masks.At(ii)
		if err != nil {
			return nil, errors.Wrapf(err, "invalid mask value at (%v)", ii)
		}
		mask := utils.FlattenInt32Slice(iMask)
		for jj := 0; jj < featureLen; jj++ {
			rprobs[jj] = feature.New(
				feature.SemanticSegmentType(),
				feature.SemanticSegmentHeight(int32(targetHeight)),
				feature.SemanticSegmentWidth(int32(targetWidth)),
				feature.SemanticSegmentIntMask(mask),
				feature.Probability(1.0),
			)
		}
		features[ii] = rprobs
	}

	return features, nil
}

func (p ImagePredictor) CreateInstanceSegmentFeatures(ctx context.Context, probabilities [][]float32, classes [][]float32, boxes [][][]float32, masks [][][][]float32, labels []string) ([]dlframework.Features, error) {
	batchSize := p.BatchSize()
	if len(probabilities) < 1 {
		return nil, errors.New("len(probabilities) < 1")
	}
	featureLen := len(probabilities[0])
	features := make([]dlframework.Features, batchSize)
	for ii := 0; ii < batchSize; ii++ {
		rprobs := make([]*dlframework.Feature, featureLen)
		for jj := 0; jj < featureLen; jj++ {
			mask := masks[ii][jj]
			masktype := "float"
			height := len(mask)
			width := len(mask[0])
			rprobs[jj] = feature.New(
				feature.InstanceSegmentType(),
				feature.InstanceSegmentXmin(boxes[ii][jj][1]),
				feature.InstanceSegmentXmax(boxes[ii][jj][3]),
				feature.InstanceSegmentYmin(boxes[ii][jj][0]),
				feature.InstanceSegmentYmax(boxes[ii][jj][2]),
				feature.InstanceSegmentIndex(int32(classes[ii][jj])),
				feature.InstanceSegmentLabel(labels[int32(classes[ii][jj])]),
				feature.InstanceSegmentMaskType(masktype),
				feature.InstanceSegmentFloatMask(flattenFloat32Slice(mask)),
				feature.InstanceSegmentHeight(int32(height)),
				feature.InstanceSegmentWidth(int32(width)),
				feature.Probability(probabilities[ii][jj]),
			)
		}
		sort.Sort(dlframework.Features(rprobs))
		features[ii] = rprobs
	}

	return features, nil
}

func (p ImagePredictor) CreateRawImageFeatures(ctx context.Context, images [][][][]float32) ([]dlframework.Features, error) {
	batchSize := p.BatchSize()
	if len(images) == 0 {
		return nil, errors.New("len(outImages) = 0")
	}
	height := len(images[0])
	width := len(images[0][0])
	channels := len(images[0][0][0])

	mean, err := p.GetMeanImage()
	if err != nil {
		return nil, err
	}
	scale, err := p.GetScale()
	if err != nil {
		return nil, err
	}

	features := make([]dlframework.Features, batchSize)
	for ii := 0; ii < batchSize; ii++ {
		curr := images[ii]
		pixels := make([]float32, width*height*channels)
		for h := 0; h < height; h++ {
			for w := 0; w < width; w++ {
				pixels[(h*width+w)*channels+0] = curr[h][w][0]*scale + mean[0]
				pixels[(h*width+w)*channels+1] = curr[h][w][1]*scale + mean[1]
				pixels[(h*width+w)*channels+2] = curr[h][w][2]*scale + mean[2]
			}
		}
		features[ii] = dlframework.Features{
			feature.New(
				feature.RawImageType(),
				feature.RawImageWidth(width),
				feature.RawImageHeight(height),
				feature.RawImageChannels(channels),
				feature.RawImageData(pixels),
			),
		}
	}

	return features, nil
}

func (p ImagePredictor) CreateImageFeatures(ctx context.Context, images [][][][]float32) ([]dlframework.Features, error) {
	batchSize := p.BatchSize()
	if len(images) < 1 {
		return nil, errors.New("len(outImages) < 1")
	}
	height := len(images[0])
	width := len(images[0][0])
	mean, err := p.GetMeanImage()
	if err != nil {
		return nil, err
	}
	scale, err := p.GetScale()
	if err != nil {
		return nil, err
	}
	features := make([]dlframework.Features, batchSize)

	for ii := 0; ii < batchSize; ii++ {
		curr := images[ii]
		img := image.NewRGBA(image.Rect(0, 0, width, height))
		for w := 0; w < width; w++ {
			for h := 0; h < height; h++ {
				R := uint8(curr[h][w][0]*scale + mean[0])
				G := uint8(curr[h][w][1]*scale + mean[1])
				B := uint8(curr[h][w][2]*scale + mean[2])
				img.Set(w, h, color.RGBA{R, G, B, 255})
			}
		}
		pp.Println(img.At(0, 0))

		buf := new(bytes.Buffer)
		err = jpeg.Encode(buf, img, nil)
		if err != nil {
			return nil, err
		}
		imgBytes := buf.Bytes()
		features[ii] = dlframework.Features{feature.New(
			feature.ImageType(),
			feature.ImageData(imgBytes),
		)}
	}

	return features, nil
}

func (p ImagePredictor) Reset(ctx context.Context) error {
	return nil
}

func (p ImagePredictor) Close() error {
	return nil
}
