package yaml

import (
	"fmt"
	"strconv"
	"vislab/sources/yaml/types"

	"gopkg.in/yaml.v3"
)

type Parser struct {
	settersMap map[string]any
}

func NewParser(configData []byte) (*Parser, error) {
	p := &Parser{}

	settersMap, err := p.parseConfig(configData)
	if err != nil {
		return nil, err
	}

	p.settersMap = settersMap

	return p, nil
}

func (p *Parser) Parse(in []byte, out *types.All) error {
	yamlMap := map[string]any{}

	if err := yaml.Unmarshal(in, &yamlMap); err != nil {
		return err
	}

	if err := p.triggerSetters(yamlMap, out); err != nil {
		return err
	}

	return nil
}

func (p *Parser) triggerSetters(in map[string]any, out *types.All) error {
	if err := iterateSettersMap(in, p.settersMap, out); err != nil {
		return err
	}

	return nil
}

func iterateSettersMap(parsedMap, settersMap map[string]any, out *types.All) error {
	for key, value := range parsedMap {
		var sMap any
		switch {
		case settersMap[key] != nil:
			sMap = settersMap[key]
		case settersMap["*"] != nil:
			sMap = settersMap["*"]
		default:
			continue
		}

		switch value := value.(type) {
		case map[string]any:
			if sMap, ok := sMap.(map[string]any); ok {
				if err := iterateSettersMap(value, sMap, out); err != nil {
					return err
				}
			}
		case []any:
			if tmpMap, ok := sMap.(map[string]any); ok && tmpMap["[]"] != nil {
				sMap = tmpMap["[]"]
			}

			if sMap, ok := sMap.([]any); ok {
				if len(sMap) == 0 {
					continue
				}

				if err := iterateSettersArray(value, sMap[0], out); err != nil {
					return err
				}
			}
		case string:
			if sMap, ok := sMap.(*Setter); ok {
				if err := sMap.Set(value, out); err != nil {
					return err
				}
			}
		case int:
			if sMap, ok := sMap.(*Setter); ok {
				if err := sMap.Set(strconv.Itoa(value), out); err != nil {
					return err
				}
			}
		case bool:
			if sMap, ok := sMap.(*Setter); ok {
				if err := sMap.Set(strconv.FormatBool(value), out); err != nil {
					return err
				}
			}
		case float64:
			if sMap, ok := sMap.(*Setter); ok {
				if err := sMap.Set(strconv.FormatFloat(value, 'f', -1, 64), out); err != nil {
					return err
				}
			}
		default:
			return fmt.Errorf("unsupported type: %T", value)
		}
	}
	return nil
}

func iterateSettersArray(array []any, settersMap any, out *types.All) error {
	for _, value := range array {
		switch value := value.(type) {
		case map[string]any:
			if sMap, ok := settersMap.(map[string]any); ok {
				if err := iterateSettersMap(value, sMap, out); err != nil {
					return err
				}
			}
		case []any:
			if sMap, ok := settersMap.([]any); ok {
				if len(sMap) == 0 {
					continue
				}

				if err := iterateSettersArray(value, sMap[0], out); err != nil {
					return err
				}
			}

		case string:
			if sMap, ok := settersMap.(*Setter); ok {
				if err := sMap.Set(value, out); err != nil {
					return err
				}
			}
		case int:
			if sMap, ok := settersMap.(*Setter); ok {
				if err := sMap.Set(strconv.Itoa(value), out); err != nil {
					return err
				}
			}
		case bool:
			if sMap, ok := settersMap.(*Setter); ok {
				if err := sMap.Set(strconv.FormatBool(value), out); err != nil {
					return err
				}
			}
		case float64:
			if sMap, ok := settersMap.(*Setter); ok {
				if err := sMap.Set(strconv.FormatFloat(value, 'f', -1, 64), out); err != nil {
					return err
				}
			}
		default:
			return fmt.Errorf("unsupported type: %T", value)
		}
	}

	return nil
}
