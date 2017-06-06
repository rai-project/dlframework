package generator

import (
	"github.com/dave/jennifer/jen"
	"github.com/rai-project/dlframework/mxnet"
)

const (
	imgIoPkg     = "github.com/anthonynsimon/bild/imgio"
	transformPkg = "github.com/anthonynsimon/bild/transform"
	mxnetPkg     = "github.com/songtianyi/go-mxnet-predictor/mxnet"
	utilsPkg     = "github.com/songtianyi/go-mxnet-predictor/utils"
)

type generator struct {
	model mxnet.ModelInformation
}

func New(model mxnet.ModelInformation) (*generator, error) {
	return &generator{model: model}, nil
}

func (g *generator) Generate() *jen.File {
	f := jen.NewFile("main")
	f.Anon(
		"github.com/anthonynsimon/bild/imgio",
		"github.com/anthonynsimon/bild/transform",
		"github.com/songtianyi/go-mxnet-predictor/mxnet",
		"github.com/songtianyi/go-mxnet-predictor/utils",
	)
	f.Func().Id("main").Params().Block(
		jen.Qual("github.com/k0kubun/pp", "Println").Call(jen.Lit("Hello, world")),
	)

	return f
}
