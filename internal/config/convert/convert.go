package convert

import (
	"encoding/json"
	"strconv"

	"gopkg.in/yaml.v3"
)

func gen(n *yaml.Node) interface{} {
	switch n.Kind {

	case yaml.MappingNode:
		nn := map[string]interface{}{}
		for i, _ := range n.Content {
			if i%2 == 0 {
				k := n.Content[i]
				nn[k.Value] = gen(n.Content[i+1])
			}
		}
		return nn

	case yaml.SequenceNode:
		nn := []interface{}{}
		for i, _ := range n.Content {
			nn = append(nn, gen(n.Content[i]))
		}
		return nn

	case yaml.ScalarNode:
		switch n.Tag {

		case "!!bool":
			return n.Value == "true"

		case "!!int":
			i, err := strconv.Atoi(n.Value)
			if err != nil {
				panic(err)
			}
			return i

		default:
			return n.Value

		}

	default:
		return ""

	}
}

func YAMLtoJSON(n yaml.Node) []byte {
	blob, err := json.Marshal(gen(&n))
	if err != nil {
		panic(err)
	}
	return blob
}
