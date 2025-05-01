package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"
)

type Out struct {
	Files []File `json:"files"`
}

type File struct {
	Name     string `json:"name"`
	Contents []byte `json:"contents"`
}

func main() {
	in := make(map[string]interface{})
	decoder := json.NewDecoder(os.Stdin)
	err := decoder.Decode(&in)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error generating JSON: %s", err)
		os.Exit(2)
	}

	buf := bytes.NewBuffer(nil)
	queries := in["queries"].([]interface{})
	for _, q := range queries {
		text := q.(map[string]interface{})["text"].(string)
		buf.WriteString(text)
		buf.WriteString("\n")
	}

	e := json.NewEncoder(os.Stdout)
	e.SetIndent("", "  ")
	e.Encode(&Out{Files: []File{{Name: "hello.txt", Contents: buf.Bytes()}}})
}
