package json

type opts struct {
	Out      string `json:"out"`
	Indent   string `json:"indent,omitempty"`
	Filename string `json:"filename,omitempty"`
}
