package config

type PythonType struct {
	Module string `json:"module" yaml:"module"`
	Name   string `json:"name" yaml:"name"`
}

func (t PythonType) IsSet() bool {
	return t.Module != "" || t.Name != ""
}

func (t PythonType) TypeString() string {
	if t.Name != "" && t.Module == "" {
		return t.Name
	}
	return t.Module + "." + t.Name
}
