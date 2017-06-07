package mxnet

import (
	"encoding/json"
	"errors"
)

func (e *Graph_NodeEntry) UnmarshalJSON(b []byte) error {
	var s []int32
	if err := json.Unmarshal(b, &s); err != nil {
		return err
	}
	if len(s) < 2 {
		return errors.New("expecting a node entry length >= 2")
	}
	e.NodeId = s[0]
	e.Index = s[1]
	if len(s) == 3 {
		e.Version = s[2]
	}
	return nil
}

func (e *Graph_NodeEntry) MarshalJSON() ([]byte, error) {
	s := []int32{
		e.NodeId,
		e.Index,
	}
	if e.GetVersion() != 0 {
		s = append(s, e.Version)
	}

	return json.Marshal(s)
}
