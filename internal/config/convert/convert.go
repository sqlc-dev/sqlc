package convert

import (
	"encoding/json"
	"fmt"
	"strconv"

	"gopkg.in/yaml.v3"
)

func gen(n *yaml.Node) (interface{}, error) {
	switch n.Kind {

	case yaml.MappingNode:
		nn := map[string]interface{}{}
		for i, _ := range n.Content {
			if i%2 == 0 {
				k := n.Content[i]
				v, err := gen(n.Content[i+1])
				if err != nil {
					return nil, err
				}
				nn[k.Value] = v
			}
		}
		return nn, nil

	case yaml.SequenceNode:
		nn := []interface{}{}
		for i, _ := range n.Content {
			v, err := gen(n.Content[i])
			if err != nil {
				return nil, err
			}
			nn = append(nn, v)
		}
		return nn, nil

	case yaml.ScalarNode:
		switch n.Tag {

		case "!!bool":
			return strconv.ParseBool(n.Value)

		case "!!float":
			return strconv.ParseFloat(n.Value, 64)

		case "!!int":
			return strconv.Atoi(n.Value)

		case "!!null":
			return nil, nil

		case "!!str":
			return n.Value, nil

		default:
			return n.Value, nil

		}

	case yaml.AliasNode:
		return gen(n.Alias)

	default:
		return nil, fmt.Errorf("unknown yaml value: %s (%s)", n.Value, n.Tag)

	}
}

func YAMLtoJSON(n yaml.Node) ([]byte, error) {
	if n.Kind == 0 {
		return []byte{}, nil
	}
	iface, err := gen(&n)
	if err != nil {
		return nil, err
	}
	blob, err := json.Marshal(iface)
	if err != nil {
		return nil, err
	}
	return blob, nil
}
