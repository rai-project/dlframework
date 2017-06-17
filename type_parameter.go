package dlframework

import (
	"errors"

	google_protobuf "github.com/gogo/protobuf/types"
)

func (param *ModelManifest_Type_Parameter) UnmarshalJSON(b []byte) error {
	return nil
}

func (param *ModelManifest_Type_Parameter) MarshalJSON() ([]byte, error) {
	return []byte{}, nil
}

func (param *ModelManifest_Type_Parameter) MarshalYAML() (interface{}, error) {
	return "", nil
}

func (param *ModelManifest_Type_Parameter) UnmarshalYAML(unmarshal func(interface{}) error) error {
	var intSlice []int
	if err := unmarshal(&intSlice); err == nil {
		this := &google_protobuf.Value_ListValue{}
		this.ListValue = &google_protobuf.ListValue{
			Values: make([]*google_protobuf.Value, len(intSlice)),
		}
		for ii, val := range intSlice {
			num := &google_protobuf.Value_NumberValue{
				NumberValue: float64(val),
			}
			this.ListValue.Values[ii] = &google_protobuf.Value{
				Kind: num,
			}
		}

		param.Value = &google_protobuf.Struct{
			Fields: map[string]*google_protobuf.Value{
				"data": {Kind: this},
			},
		}
		return nil
	}
	var floatSlice []float64
	if err := unmarshal(&floatSlice); err == nil {
		this := &google_protobuf.Value_ListValue{}
		this.ListValue = &google_protobuf.ListValue{
			Values: make([]*google_protobuf.Value, len(floatSlice)),
		}
		for ii, val := range floatSlice {
			num := &google_protobuf.Value_NumberValue{
				NumberValue: float64(val),
			}
			this.ListValue.Values[ii] = &google_protobuf.Value{
				Kind: num,
			}
		}

		param.Value = &google_protobuf.Struct{
			Fields: map[string]*google_protobuf.Value{
				"data": {Kind: this},
			},
		}
		return nil
	}
	var strSlice []string
	if err := unmarshal(&strSlice); err == nil {
		this := &google_protobuf.Value_ListValue{}
		this.ListValue = &google_protobuf.ListValue{
			Values: make([]*google_protobuf.Value, len(strSlice)),
		}
		for ii, val := range strSlice {
			str := &google_protobuf.Value_StringValue{
				StringValue: val,
			}
			this.ListValue.Values[ii] = &google_protobuf.Value{
				Kind: str,
			}
		}

		param.Value = &google_protobuf.Struct{
			Fields: map[string]*google_protobuf.Value{
				"data": {Kind: this},
			},
		}
		return nil
	}
	var byteSlice []byte
	if err := unmarshal(&byteSlice); err == nil {
		str := string(byteSlice)
		this := &google_protobuf.Value_StringValue{
			StringValue: str,
		}
		param.Value = &google_protobuf.Struct{
			Fields: map[string]*google_protobuf.Value{
				"data": {Kind: this},
			},
		}
		return nil
	}
	var str string
	if err := unmarshal(&str); err == nil {
		this := &google_protobuf.Value_StringValue{
			StringValue: str,
		}
		param.Value = &google_protobuf.Struct{
			Fields: map[string]*google_protobuf.Value{
				"data": {Kind: this},
			},
		}
		return nil
	}
	var b bool
	if err := unmarshal(&b); err == nil {
		this := &google_protobuf.Value_BoolValue{
			BoolValue: b,
		}
		param.Value = &google_protobuf.Struct{
			Fields: map[string]*google_protobuf.Value{
				"data": {Kind: this},
			},
		}
		return nil
	}
	return errors.New("unable to unmarshal model type parameter")
}
