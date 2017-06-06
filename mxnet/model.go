package mxnet

type Model struct {
	Name       string   `json:"name,omitempty" yaml:"name,omitempty"`
	Framework  string   `json:"description,omitempty" yaml:"description,omitempty"`
	Version    float64  `json:"version,omitempty" yaml:"version,omitempty"`
	Type       string   `json:"type,omitempty" yaml:"type,omitempty"`
	Dataset    string   `json:"dataset,omitempty" yaml:"dataset,omitempty"`
	Model      string   `json:"model,omitempty" yaml:"model,omitempty"`
	Weights    string   `json:"weights,omitempty" yaml:"weights,omitempty"`
	References []string `json:"references,omitempty" yaml:"references,omitempty"`
}
