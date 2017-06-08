package main

import (
	"bufio"
	"flag"
	"fmt"
	"image"
	"io"
	"io/ioutil"
	"log"
	"math"
	"net/http"
	"os"
	"path/filepath"
	"sort"

	"github.com/anthonynsimon/bild/imgio"
	"github.com/anthonynsimon/bild/transform"
	"github.com/songtianyi/go-mxnet-predictor/mxnet"
	"github.com/songtianyi/go-mxnet-predictor/utils"
)

const path = "/Users/chengli/Downloads"

var (
	caffenet    = []uint32{1, 3, 224, 224}
	rn1015k500  = []uint32{1, 3, 224, 224}
	vgg19       = []uint32{1, 3, 224, 224}
	inceptionbn = []uint32{1, 3, 224, 224}
)

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

func preprocessImage(src image.Image) ([]float32, error) {

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

func getImageWithURL(url string) (err error) {
	response, e := http.Get(url)
	if e != nil {
		return err
	}
	defer response.Body.Close()

	file, err := os.Create("/tmp/tmp.jpg")
	if err != nil {
		return err
	}
	_, err = io.Copy(file, response.Body)
	if err != nil {
		log.Fatal(err)
	}
	file.Close()

	return nil
}

func getModelShape(prefix string) []uint32 {
	switch prefix {
	case "caffenet":
		return caffenet
	case "RN101-5k500":
		return rn1015k500
	case "vgg19":
		return vgg19
	case "Inception-BN":
		return inceptionbn
	default:
		return rn1015k500
	}
}

func Predict(dirpath string, prefix string, epoch string, labelfile string, image string) (prob float32, output string, err error) {
	symbol, err := ioutil.ReadFile(filepath.Join(dirpath, prefix+"-symbol.json"))
	if err != nil {
		return 0, "", err
	}
	params, err := ioutil.ReadFile(filepath.Join(dirpath, prefix+"-"+epoch+".params"))
	if err != nil {
		return 0, "", err
	}

	var labels []string
	f, _ := os.Open(filepath.Join(path, labelfile))
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := scanner.Text()
		labels = append(labels, line)
	}

	inputShape := getModelShape(prefix)
	p, err := mxnet.CreatePredictor(symbol,
		params,
		mxnet.Device{mxnet.CPU_DEVICE, 0},
		[]mxnet.InputNode{{Key: "data", Shape: inputShape}},
	)
	if err != nil {
		return 0, "", err
	}
	defer p.Free()

	img, err := imgio.Open(filepath.Join(path, image))
	if err != nil {
		panic(err)
	}

	resized := transform.Resize(img, int(inputShape[2]), int(inputShape[3]), transform.Linear)
	res, err := preprocessImage(resized)
	if err != nil {
		return 0, "", err
	}

	if err := p.SetInput("data", res); err != nil {
		return 0, "", err
	}

	if err := p.Forward(); err != nil {
		return 0, "", err
	}

	probs, err := p.GetOutput(0)
	if err != nil {
		panic(err)
	}
	idxs := make([]int, len(probs))
	for i := range probs {
		idxs[i] = i
	}

	out := utils.ArgSort{Args: probs, Idxs: idxs}
	sort.Sort(out)

	fmt.Println("result:")
	fmt.Println(out.Args[0])
	fmt.Println(labels[out.Idxs[0]])

	return out.Args[0], labels[out.Idxs[0]], nil
}

func main() {
	flag.Parse()
	args := flag.Args()
	if len(args) < 4 {
		fmt.Println("Input error")
		os.Exit(1)
	}
	// example RN101-5k500 0012 grids.txt tokyo-tower.jpg
	Predict(path, args[0], args[1], args[2], args[3])
}
