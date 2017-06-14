package dlframework

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
	return nil
}
