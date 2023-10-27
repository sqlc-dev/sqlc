package config

import (
	"encoding/json"
)

type GoType struct {
	Path    string `json:"import" yaml:"import"`
	Package string `json:"package" yaml:"package"`
	Name    string `json:"type" yaml:"type"`
	Pointer bool   `json:"pointer" yaml:"pointer"`
	Slice   bool   `json:"slice" yaml:"slice"`
	Spec    string
	BuiltIn bool
}

func (o *GoType) UnmarshalJSON(data []byte) error {
	var spec string
	if err := json.Unmarshal(data, &spec); err == nil {
		*o = GoType{Spec: spec}
		return nil
	}
	type alias GoType
	var a alias
	if err := json.Unmarshal(data, &a); err != nil {
		return err
	}
	*o = GoType(a)
	return nil
}

func (o *GoType) UnmarshalYAML(unmarshal func(interface{}) error) error {
	var spec string
	if err := unmarshal(&spec); err == nil {
		*o = GoType{Spec: spec}
		return nil
	}
	type alias GoType
	var a alias
	if err := unmarshal(&a); err != nil {
		return err
	}
	*o = GoType(a)
	return nil
}
