package predictor

import (
	"bytes"
	goimage "image"
	"image/color"
	"image/jpeg"
	"path/filepath"
	"sort"

	"github.com/k0kubun/pp"

	"context"

	"github.com/pkg/errors"
	"github.com/rai-project/dlframework"
	"github.com/rai-project/dlframework/framework/feature"
	"github.com/rai-project/image"
	"github.com/rai-project/image/types"
	imageTypes "github.com/rai-project/image/types"
	yaml "gopkg.in/yaml.v2"
)

type PreprocessOptions struct {
	Context         context.Context
	ElementType     string
	MeanImage       []float32
	Dims            []int
	MaxDimension    *int
	KeepAspectRatio *bool
	Scale           float32
	ColorMode       types.Mode
	Layout          image.Layout
}

type ImagePredictor struct {
	Base
	Metadata map[string]interface{}
}

func (p ImagePredictor) GetImageDimensions() ([]int, error) {
	model := p.Model
	modelInputs := model.GetInputs()
	typeParameters := modelInputs[0].GetParameters()
	if typeParameters == nil {
		return nil, errors.New("invalid type parameters")
	}
	pdims, ok := typeParameters["dimensions"]
	if !ok {
		log.Debug("arbitrary image dimensions")
		return nil, nil
	}
	pdimsVal := pdims.GetValue()
	if pdimsVal == "" {
		return nil, errors.New("invalid image dimensions")
	}

	var dims []int
	if err := yaml.Unmarshal([]byte(pdimsVal), &dims); err != nil {
		return nil, errors.Errorf("unable to get image dimensions %v as an integer slice", pdimsVal)
	}
	if len(dims) == 1 {
		dims = []int{dims[0], dims[0], dims[0]}
	}
	if len(dims) > 3 {
		return nil, errors.Errorf("expecting a dimensions size of 1 or 3, but got %v. do not put the batch size in the input dimensions.", len(dims))
	}

	return dims, nil
}

func (p ImagePredictor) GetMeanPath() string {
	model := p.Model
	return cleanString(filepath.Join(p.WorkDir, model.GetName()+".mean"))
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

func (p ImagePredictor) GetLayout(defaultLayout image.Layout) image.Layout {
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
		return image.CHWLayout
	case "HWC":
		return image.HWCLayout
	default:
		log.Error("invalid image mode specified " + val)
		return image.InvalidLayout
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

func (p ImagePredictor) GetPreprocessOptions(ctx context.Context) (PreprocessOptions, error) {
	mean, err := p.GetMeanImage()
	if err != nil {
		return PreprocessOptions{}, err
	}
	scale, err := p.GetScale()
	if err != nil {
		return PreprocessOptions{}, err
	}

	imageDims, err := p.GetImageDimensions()
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
		Context:         ctx,
		ElementType:     p.Model.GetElementType(),
		MeanImage:       mean,
		Scale:           scale,
		Dims:            imageDims,
		MaxDimension:    maxDim,
		KeepAspectRatio: keepAspectRatio,
		ColorMode:       p.GetColorMode(imageTypes.RGBMode),
		Layout:          p.GetLayout(image.HWCLayout),
	}

	return preprocOpts, nil
}

func (p ImagePredictor) CreateClassificationFeatures(ctx context.Context, probabilities [][]float32, labels []string) ([]dlframework.Features, error) {
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

func (p ImagePredictor) CreateBoundingBoxFeatures(ctx context.Context, probabilities [][]float32, classes [][]float32, boxes [][][]float32, labels []string) ([]dlframework.Features, error) {
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

func (p ImagePredictor) CreateSemanticSegmentFeatures(ctx context.Context, masks [][][]int64, labels []string) ([]dlframework.Features, error) {
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
	channels := 3

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
		pixels := make([]uint8, width*height*channels)
		for h := 0; h < height; h++ {
			for w := 0; w < width; w++ {
				R := uint8(curr[h][w][0]*scale + mean[0])
				G := uint8(curr[h][w][1]*scale + mean[1])
				B := uint8(curr[h][w][2]*scale + mean[2])
				pixels[(h*width+w)*channels+0] = R
				pixels[(h*width+w)*channels+1] = G
				pixels[(h*width+w)*channels+2] = B
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
		img := goimage.NewRGBA(goimage.Rect(0, 0, width, height))
		var R, G, B uint8
		for w := 0; w < width; w++ {
			for h := 0; h < height; h++ {
				R, G, B = uint8(curr[h][w][0]*scale+mean[0]), uint8(curr[h][w][1]*scale+mean[1]), uint8(curr[h][w][2]*scale+mean[2])
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
