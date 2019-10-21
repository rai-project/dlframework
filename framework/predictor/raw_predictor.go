package predictor

import (
	"bufio"
	"context"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/pkg/errors"
	"github.com/rai-project/dlframework"
	"github.com/rai-project/downloadmanager"
)

const rawFileDelimiter = ","

type RawPredictor struct {
	Base
	Metadata map[string]interface{}
}

func (p RawPredictor) GetInputParams(name string) ([]string, error) {
	model := p.Model
	modelInputs := model.GetInputs()
	params := []string{}
	for _, input := range modelInputs {
		typeParameters := input.GetParameters()
		param, err := p.GetTypeParameter(typeParameters, name)
		if err != nil {
			continue
		}
		params = append(params, param)
	}
	return params, nil
}

func (p RawPredictor) GetInputParamsByIdx(idx int, name string) (string, error) {
	model := p.Model
	modelInputs := model.GetInputs()
	if idx > len(modelInputs) {
		return "", errors.New("idx is larger than the number of inputs")
	}
	input := modelInputs[idx]
	typeParameters := input.GetParameters()
	return p.GetTypeParameter(typeParameters, name)
}

func (p RawPredictor) GetInputURLs() ([]string, error) {
	model := p.Model
	modelInputs := model.GetInputs()
	urls := []string{}
	for _, input := range modelInputs {
		typeParameters := input.GetParameters()
		url, err := p.GetTypeParameter(typeParameters, "url")
		if err != nil {
			continue
		}
		url = strings.TrimRight(url, "/")
		url = strings.TrimSpace(url)
		urls = append(urls, url)
	}
	return urls, nil
}

func (p RawPredictor) GetInputDataByIdx(idx int) ([]interface{}, error) {
	urls, err := p.GetInputURLs()
	if err != nil {
		return nil, err
	}
	inputLayer, err := p.GetInputParamsByIdx(idx, "input_layer")
	if err != nil {
		return nil, err
	}
	inputType, err := p.GetInputParamsByIdx(idx, "input_type")
	if err != nil {
		return nil, err
	}
	elementType, err := p.GetInputParamsByIdx(idx, "element_type")
	if err != nil {
		return nil, err
	}

	url := urls[idx]
	targetPath := filepath.Join(p.WorkDir, inputLayer)
	if url != "" {
		_, _, err := downloadmanager.DownloadFile(
			url,
			targetPath,
		)
		if err != nil {
			return nil, err
		}
	}
	f, err := os.Open(targetPath)
	if err != nil {
		return nil, errors.Wrapf(err, "cannot read %s", targetPath)
	}

	batchSize := p.BatchSize()
	cnt := 0

	var data []interface{}
	switch inputType {
	case "scalar":
		switch elementType {
		case "int8":
			scanner := bufio.NewScanner(f)
			for scanner.Scan() {
				if cnt == batchSize {
					break
				}
				line := scanner.Text()
				v, err := strconv.ParseInt(line, 10, 8)
				if err != nil {
					return nil, err
				}
				data = append(data, int8(v))
				cnt++
			}
		case "int16":
			scanner := bufio.NewScanner(f)
			for scanner.Scan() {
				if cnt == batchSize {
					break
				}
				line := scanner.Text()
				v, err := strconv.ParseInt(line, 10, 16)
				if err != nil {
					return nil, err
				}
				data = append(data, int16(v))
				cnt++
			}
		case "int32":
			scanner := bufio.NewScanner(f)
			for scanner.Scan() {
				if cnt == batchSize {
					break
				}
				line := scanner.Text()
				v, err := strconv.ParseInt(line, 10, 32)
				if err != nil {
					return nil, err
				}
				data = append(data, int32(v))
				cnt++
			}
		case "int64":
			scanner := bufio.NewScanner(f)
			for scanner.Scan() {
				if cnt == batchSize {
					break
				}
				line := scanner.Text()
				v, err := strconv.ParseInt(line, 10, 64)
				if err != nil {
					return nil, err
				}
				data = append(data, int32(v))
				cnt++
			}
		case "float32":
			scanner := bufio.NewScanner(f)
			for scanner.Scan() {
				if cnt == batchSize {
					break
				}
				line := scanner.Text()
				v, err := strconv.ParseFloat(line, 32)
				if err != nil {
					return nil, err
				}
				data = append(data, float32(v))
				cnt++
			}
		case "float64":
			scanner := bufio.NewScanner(f)
			for scanner.Scan() {
				if cnt == batchSize {
					break
				}
				line := scanner.Text()
				v, err := strconv.ParseFloat(line, 64)
				if err != nil {
					return nil, err
				}
				data = append(data, float64(v))
				cnt++
			}
		default:
			return nil, errors.Errorf("the scalar element type=%s is not valid", elementType)
		}

	case "slice":
		switch elementType {
		case "int8":
			scanner := bufio.NewScanner(f)
			for scanner.Scan() {
				if cnt == batchSize {
					break
				}
				dd := []int8{}
				line := scanner.Text()
				s := strings.Split(line, rawFileDelimiter)
				for _, ss := range s {
					if ss == "" {
						continue
					}

					v, err := strconv.ParseInt(ss, 10, 8)
					if err != nil {
						return nil, err
					}
					dd = append(dd, int8(v))
				}
				data = append(data, dd)
				cnt++
			}
		case "int16":
			scanner := bufio.NewScanner(f)
			for scanner.Scan() {
				if cnt == batchSize {
					break
				}
				dd := []int16{}
				line := scanner.Text()
				s := strings.Split(line, rawFileDelimiter)
				for _, ss := range s {
					if ss == "" {
						continue
					}

					v, err := strconv.ParseInt(ss, 10, 16)
					if err != nil {
						return nil, err
					}
					dd = append(dd, int16(v))
				}
				data = append(data, dd)
				cnt++
			}
		case "int32":
			scanner := bufio.NewScanner(f)
			for scanner.Scan() {
				if cnt == batchSize {
					break
				}
				dd := []int32{}
				line := scanner.Text()
				s := strings.Split(line, rawFileDelimiter)
				for _, ss := range s {
					if ss == "" {
						continue
					}

					v, err := strconv.ParseInt(ss, 10, 32)
					if err != nil {
						return nil, err
					}
					dd = append(dd, int32(v))
				}
				data = append(data, dd)
				cnt++
			}
		case "int64":
			scanner := bufio.NewScanner(f)
			for scanner.Scan() {
				if cnt == batchSize {
					break
				}
				dd := []int64{}
				line := scanner.Text()
				s := strings.Split(line, rawFileDelimiter)
				for _, ss := range s {
					if ss == "" {
						continue
					}

					v, err := strconv.ParseInt(ss, 10, 64)
					if err != nil {
						return nil, err
					}
					dd = append(dd, int64(v))
				}
				data = append(data, dd)
				cnt++
			}
		case "float32":
			scanner := bufio.NewScanner(f)
			for scanner.Scan() {
				if cnt == batchSize {
					break
				}
				dd := []float32{}
				line := scanner.Text()
				s := strings.Split(line, rawFileDelimiter)
				for _, ss := range s {
					if ss == "" {
						continue
					}

					v, err := strconv.ParseFloat(ss, 32)
					if err != nil {
						return nil, err
					}
					dd = append(dd, float32(v))
				}
				data = append(data, dd)
				cnt++
			}
		case "float64":
			scanner := bufio.NewScanner(f)
			for scanner.Scan() {
				if cnt == batchSize {
					break
				}
				dd := []float64{}
				line := scanner.Text()
				s := strings.Split(line, rawFileDelimiter)
				for _, ss := range s {
					if ss == "" {
						continue
					}
					v, err := strconv.ParseFloat(ss, 32)
					if err != nil {
						return nil, err
					}
					dd = append(dd, float64(v))
				}
				data = append(data, dd)
				cnt++
			}
		default:
			return nil, errors.Errorf("the slice element type=%s is not valid", elementType)
		}

	default:
		return nil, errors.New("input type not supported")
	}
	return data, nil

}

func (p RawPredictor) GetOutputLayerName(layer string) (string, error) {
	model := p.Model
	modelOutput := model.GetOutput()

	typeParameters := modelOutput.GetParameters()

	param, err := p.GetTypeParameter(typeParameters, layer)
	if err != nil {
		return "", err
	}
	return param, nil
}

func (p RawPredictor) ReadPredictedFeatures(ctx context.Context) ([]dlframework.Features, error) {
	return []dlframework.Features{}, nil
}

func (p RawPredictor) Reset(ctx context.Context) error {
	return nil
}

func (p RawPredictor) Close() error {
	return nil
}
