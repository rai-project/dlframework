package graph

import (
	"encoding/json"
	"strings"

	"github.com/awalterschulze/gographviz"
	"github.com/pkg/errors"
	"gopkg.in/fatih/set.v0"
)

type NodeEntry struct {
	NodeId  int  `json:"node_id"`
	Index   int  `json:"index"`
	Version *int `json:"version,omitempty"`
}

type Node struct {
	Op                  string            `json:"op"`
	Param               map[string]string `json:"param"`
	Name                string            `json:"name"`
	Inputs              []NodeEntry       `json:"inputs"`
	BackwardSourceID    int               `json:"backward_source_id"`
	ControlDependencies []int             `json:"control_deps,omitempty"`
}

type Graph struct {
	Nodes          []Node                 `json:"nodes"`
	ArgNodes       []int                  `json:"arg_nodes"`
	NodeRowPointer []int                  `json:"node_row_ptr,omitempty"`
	Heads          []NodeEntry            `json:"heads"`
	Attributes     map[string]interface{} `json:"attrs,omitempty"`
}

func (e *NodeEntry) UnmarshalJSON(b []byte) error {
	var s []int
	if err := json.Unmarshal(b, &s); err != nil {
		return err
	}
	if len(s) < 2 {
		return errors.New("expecting a node entry length >= 2")
	}
	e.NodeId = s[0]
	e.Index = s[1]
	if len(s) == 3 {
		*e.Version = s[2]
	}
	return nil
}

func (e *NodeEntry) MarshalJSON() ([]byte, error) {
	s := []int{
		e.NodeId,
		e.Index,
	}
	if e.Version != nil {
		s = append(s, *e.Version)
	}

	return json.Marshal(s)
}

func (g *Graph) ToDotGraph() (*gographviz.Graph, error) {
	// color map
	var fillcolors = []string{
		"#8dd3c7",
		"#fb8072",
		"#ffffb3",
		"#bebada",
		"#80b1d3",
		"#fdb462",
		"#b3de69",
		"#fccde5",
	}
	var edgecolors = []string{
		"#245b51",
		"#941305",
		"#999900",
		"#3b3564",
		"#275372",
		"#975102",
		"#597d1c",
		"#90094e",
	}

	makeDefaultAttributes := func() map[string]string {
		return map[string]string{
			"shape":     "box",
			"fixedsize": "true",
			"width":     "1.3",
			"height":    "0.8034",
			"style":     "filled",
		}
	}

	isLikeWeight := func(name string) bool {
		if strings.HasSuffix(name, "_weight") {
			return true
		}
		if strings.HasSuffix(name, "_bias") {
			return true
		}
		if strings.HasSuffix(name, "_beta") ||
			strings.HasSuffix(name, "_gamma") ||
			strings.HasSuffix(name, "_moving_var") ||
			strings.HasSuffix(name, "_moving_mean") {
			return true
		}
		return false
	}

	hideWeights := true // TODO: should be an option
	drawShape := true   // TODO: should be an option

	dg := gographviz.NewGraph()

	hiddenNodes := set.NewNonTS()

	// make nodes
	for _, node := range g.Nodes {
		op := node.Op
		name := node.Name
		attrs := makeDefaultAttributes()
		label := op

		switch op {
		case "null":
			if isLikeWeight(name) {
				if hideWeights {
					hiddenNodes.Add(name)
				}
			}
			attrs["shape"] = "oval"
			attrs["fillcolor"] = fillcolors[0]
			label = name
		case "Convolution":
		//...
		case "FullyConnected":
			//...
		}

		attrs["label"] = label
		dg.AddNode("G", name, attrs)
	}

	// make edges
	for _, node := range g.Nodes {
		op := node.Op
		name := node.Name
		if op == "null" {
			continue
		}
		inputs := node.Inputs
		for _, item := range inputs {
			inputNode := g.Nodes[item.NodeId]
			inputName := inputNode.Name
			if hiddenNodes.Has(inputName) {
				continue
			}
			attrs := map[string]string{
				"dir":       "back",
				"arrowtail": "open",
			}
			if drawShape {
				// ...
				_ = inputNode
				_ = attrs
			}
			dg.AddEdge(name, inputName, true, attrs)
		}

	}

	return nil, errors.New("not implemented")
}
