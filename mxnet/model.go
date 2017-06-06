package mxnet

type Model struct {
	Name       string   `json:"name,omitempty" yaml:"name,omitempty"`
	Framework  string   `json:"description,omitempty" yaml:"description,omitempty"`
	Version    string   `json:"version,omitempty" yaml:"version,omitempty"`
	Type       string   `json:"type,omitempty" yaml:"type,omitempty"`
	DatasetURL string   `json:"dataset_url,omitempty" yaml:"dataset_url,omitempty"`
	GraphURL   string   `json:"graph_url,omitempty" yaml:"graph_url,omitempty"`
	WeightsURL string   `json:"weights_url,omitempty" yaml:"weights_url,omitempty"`
	References []string `json:"references,omitempty" yaml:"references,omitempty"`
}
