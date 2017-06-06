package generator

import (
	"github.com/dave/jennifer/jen"
	"github.com/rai-project/dlframework/mxnet"
)

type generator struct {
	model mxnet.Model
}

func New(model mxnet.Model) (*generator, error) {
	return &generator{model: model}, nil
}

func (g *generator) Generate() *jen.File {
	f := jen.NewFile("main")
	f.Func().Id("main").Params().Block(
		jen.Qual("fmt", "Println").Call(jen.Lit("Hello, world")),
	)

	return f
}
