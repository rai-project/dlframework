// Copyright 2016 go-mxnet-predictor Authors. All Rights Reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package main

import (
	"bufio"
	"fmt"
	"image"
	"io/ioutil"
	"math"
	"os"
	"path/filepath"
	"sort"

	"github.com/anthonynsimon/bild/imgio"
	"github.com/anthonynsimon/bild/transform"
	"github.com/songtianyi/go-mxnet-predictor/mxnet"
	"github.com/songtianyi/go-mxnet-predictor/utils"
)

const path = "/Users/chengli/Downloads"

func distance(p1 []float64, p2 []float64) float64 {
	R := 6371.0
	lat1, lng1, lat2, lng2 := p1[0], p1[1], p2[0], p2[1]
	dlat := lat2 - lat1
	dlng := lng2 - lng1
	a := math.Pow(math.Sin(dlat*0.5), 2) + math.Cos(lat1)*math.Cos(lat2)*math.Pow(math.Sin(dlng*0.5), 2)
	return 2 * R * math.Asin(math.Sqrt(a))
}

// convert go Image to 1-dim array
func imageTo1DArray(src image.Image) ([]float32, error) {

	if src == nil {
		return nil, fmt.Errorf("src image nil")
	}

	b := src.Bounds()
	h := b.Max.Y - b.Min.Y // image height
	w := b.Max.X - b.Min.X // image width

	res := make([]float32, 3*h*w)
	for y := 0; y < h; y++ {
		for x := 0; x < w; x++ {
			r, g, b, _ := src.At(x+b.Min.X, y+b.Min.Y).RGBA()
			res[y*w+x] = float32(r >> 8)
			res[w*h+y*w+x] = float32(g >> 8)
			res[2*w*h+y*w+x] = float32(b >> 8)
		}
	}
	return res, nil
}

func imageTo1DArraySubMean(src image.Image) ([]float32, error) {

	if src == nil {
		return nil, fmt.Errorf("src image nil")
	}

	b := src.Bounds()
	h := b.Max.Y - b.Min.Y // image height
	w := b.Max.X - b.Min.X // image width

	res := make([]float32, 3*h*w)
	for y := 0; y < h; y++ {
		for x := 0; x < w; x++ {
			r, g, b, _ := src.At(x+b.Min.X, y+b.Min.Y).RGBA()
			res[y*w+x] = float32(r>>8) - 123.68
			res[w*h+y*w+x] = float32(g>>8) - 116.779
			res[2*w*h+y*w+x] = float32(b>>8) - 103.939
		}
	}
	return res, nil
}
func main() {
	// load model
	symbol, err := ioutil.ReadFile(filepath.Join(path, "RN101-5k500-symbol.json"))
	if err != nil {
		panic(err)
	}
	params, err := ioutil.ReadFile(filepath.Join(path, "RN101-5k500-0012.params"))
	if err != nil {
		panic(err)
	}

	var labels []string
	f, _ := os.Open(filepath.Join(path, "grids.txt"))
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := scanner.Text()
		labels = append(labels, line)
	}

	// create predictor
	p, err := mxnet.CreatePredictor(symbol,
		params,
		mxnet.Device{mxnet.CPU_DEVICE, 0},
		[]mxnet.InputNode{{Key: "data", Shape: []uint32{1, 3, 224, 224}}},
	)
	if err != nil {
		panic(err)
	}
	defer p.Free()

	// load test image for predction
	img, err := imgio.Open(filepath.Join(path, "tokyo-tower.jpg"))
	if err != nil {
		panic(err)
	}
	// preprocess
	resized := transform.Resize(img, 224, 224, transform.Linear)
	res, err := imageTo1DArraySubMean(resized)
	if err != nil {
		panic(err)
	}

	// set input
	if err := p.SetInput("data", res); err != nil {
		panic(err)
	}
	// do predict
	if err := p.Forward(); err != nil {
		panic(err)
	}
	// get predict result
	data, err := p.GetOutput(0)
	if err != nil {
		panic(err)
	}
	idxs := make([]int, len(data))
	for i := range data {
		idxs[i] = i
	}
	as := utils.ArgSort{Args: data, Idxs: idxs}
	sort.Sort(as)
	fmt.Println("result:")
	fmt.Println(as.Args[0])
	fmt.Println(labels[as.Idxs[0]])
}
