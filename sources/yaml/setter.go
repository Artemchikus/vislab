package yaml

import (
	"fmt"
	"log"
	"slices"
	"strconv"
	"strings"
	"vislab/sources/yaml/types"

	"gopkg.in/yaml.v3"
)

var alphabet = []byte("abcdefghijklmnopqrstuvwxyz_.")

type Setter struct {
	set    func(s string, all *types.All) error
	weight int
}

func (s *Setter) Set(value string, all *types.All) error {
	return s.set(value, all)
}

func (s *Setter) Weight() int {
	return s.weight
}

func (p *Parser) parseConfig(configData []byte) (map[string]any, error) {
	settersMap := map[string]any{}

	escapedConfigData := strings.ReplaceAll(string(configData), " {{", " \"{{")
	escapedConfigData = strings.ReplaceAll(escapedConfigData, "}}\n", "}}\"\n")
	escapedConfigData = strings.ReplaceAll(escapedConfigData, " *:", " \"*\":")

	if err := yaml.Unmarshal([]byte(escapedConfigData), &settersMap); err != nil {
		panic(err)
	}

	if err := setupSetters(settersMap); err != nil {
		return nil, err
	}

	return settersMap, nil
}

func getSetObjFunc(objPath string) (func(string, *types.All) error, error) {
	parts := strings.Split(objPath, ".")
	if len(parts) < 3 {
		return nil, fmt.Errorf("invalid obj path %s", objPath)
	}

	switch parts[1] {
	case "service":
		return getSetSvcFunc(parts[2:])
	case "postgresql":
		return getSetPqFunc(parts[2:])
	case "kafka":
		return getSetKafkaFunc(parts[2:])
	case "redis":
		return getSetRedisFunc(parts[2:])
	case "rabbitmq":
		return getSetRabbitFunc(parts[2:])
	case "other_service":
		return getSetOtherSvcFunc(parts[2:])
	default:
		return nil, fmt.Errorf("invalid obj path %s", objPath)
	}
}

func setupSetters(configMap map[string]any) error {
	if err := iterateConfigKeys(configMap); err != nil {
		return err
	}

	return nil
}

func iterateConfigKeys(configMap map[string]any) error {
	for key, value := range configMap {
		switch value := value.(type) {
		case map[string]any:
			if err := iterateConfigKeys(value); err != nil {
				return err
			}
		case []any:
			if err := iterateConfigArray(value); err != nil {
				return err
			}
		default:
			val := fmt.Sprintf("%s", value)

			setter, err := parseSetter(val)
			if err != nil {
				return err
			}

			configMap[key] = setter
		}
	}

	return nil
}

func iterateConfigArray(configArray []any) error {
	for key, value := range configArray {
		switch value := value.(type) {
		case map[string]any:
			if err := iterateConfigKeys(value); err != nil {
				return err
			}
		case []any:
			if err := iterateConfigArray(value); err != nil {
				return err
			}
		default:
			val := fmt.Sprintf("%s", value)

			setter, err := parseSetter(val)
			if err != nil {
				return err
			}

			configArray[key] = setter
		}
	}

	return nil
}

func parseSetter(value string) (*Setter, error) {
	parts := parseParts(value)

	funcs := []func(string, *types.All) error{}
	parsedSetter := &Setter{
		weight: 0,
	}

	for _, part := range parts {
		part = strings.Trim(part, " ")

		switch {
		case strings.HasPrefix(part, "weight "):
			weight, err := parseWeight(part)
			if err != nil {
				return nil, err
			}

			parsedSetter.weight = weight
		case strings.HasPrefix(part, "if "):
			objPath, preSet, condition, err := parseIf(part)
			if err != nil {
				return nil, err
			}

			setF, err := getSetObjFunc(objPath)
			if err != nil {
				return nil, err
			}

			f := func(s string, all *types.All) error {
				b, err := strconv.ParseBool(s)
				if err != nil {
					return fmt.Errorf("could not parse given value %w", err)
				}

				if b == condition {
					return setF(preSet, all)
				}

				return nil
			}
			funcs = append(funcs, f)
		case strings.HasPrefix(part, "parse "):
			parseParts, separators, err := parseParse(part)
			if err != nil {
				return nil, err
			}

			setFs := []func(string, *types.All) error{}

			for _, objPath := range parseParts {
				setF, err := getSetObjFunc(objPath)
				if err != nil {
					return nil, err
				}

				setFs = append(setFs, setF)
			}

			f := func(s string, all *types.All) error {
				parsedValues, err := parseStrWithSeparators(s, separators)
				if err != nil {
					return err
				}

				for i, parsedValue := range parsedValues {
					if i >= len(setFs) {
						return fmt.Errorf("too many values parsed %v", parsedValues)
					}

					if err := setFs[i](parsedValue, all); err != nil {
						return err
					}
				}
				return nil
			}

			funcs = append(funcs, f)
		case strings.Contains(part, " = "):
			objPath, preSet, err := parsePreset(part)
			if err != nil {
				return nil, err
			}

			setF, err := getSetObjFunc(objPath)
			if err != nil {
				return nil, err
			}

			f := func(s string, all *types.All) error {
				return setF(preSet, all)
			}

			funcs = append(funcs, f)
		default:
			setF, err := getSetObjFunc(part)
			if err != nil {
				return nil, err
			}

			funcs = append(funcs, setF)
		}
	}

	parsedSetter.set = func(s string, all *types.All) error {
		for _, f := range funcs {
			if err := f(s, all); err != nil {
				return err
			}
		}

		return nil
	}

	return parsedSetter, nil
}

func parseParts(value string) []string {
	value = strings.TrimSuffix(value, " }}")
	value = strings.TrimPrefix(value, "{{ ")
	parts := strings.Split(value, " | ")

	return parts
}

func parseWeight(part string) (int, error) {
	trimmedPart := strings.TrimPrefix(part, "weight ")

	w, err := strconv.Atoi(trimmedPart)
	if err != nil {
		return 0, fmt.Errorf("could not parse weight %w", err)
	}

	return w, nil
}

func parseIf(part string) (string, string, bool, error) {
	trimmedPart := strings.TrimPrefix(part, "if ")

	ifParts := strings.SplitN(trimmedPart, " ", 2)
	if len(ifParts) != 2 {
		return "", "", false, fmt.Errorf("invalid if part %s, should be 'if <true|false> <key> = <value>'", part)
	}

	condition, err := strconv.ParseBool(ifParts[0])
	if err != nil {
		return "", "", false, fmt.Errorf("could not parse condition %w", err)
	}

	vParts := strings.SplitN(ifParts[1], " = ", 2)
	if len(vParts) != 2 {
		return "", "", false, fmt.Errorf("invalid if part %s, should be 'if <true|false> <key> = <value>'", part)
	}

	objPath, preSet := vParts[0], vParts[1]

	return objPath, preSet, condition, nil
}

func parseParse(part string) ([]string, []string, error) {
	trimmedPart := strings.TrimPrefix(part, "parse ")

	parseParts := []string{}
	var parsePart string
	separators := []string{}
	var separator string

	for i := 0; i < len(trimmedPart); i++ {
		if trimmedPart[i] == '.' {
			if separator != "" {
				separators = append(separators, separator)
				separator = ""
			}

			var j int
			for j = i; j < len(trimmedPart); j++ {
				if !slices.Contains(alphabet, trimmedPart[j]) {
					break
				}

				parsePart += string(trimmedPart[j])
			}

			i = j
			parseParts = append(parseParts, parsePart)
			parsePart = ""
		}

		if i < len(trimmedPart) {
			separator += string(trimmedPart[i])
		}
	}

	if separator != "" {
		separators = append(separators, separator)
	}

	return parseParts, separators, nil
}

func parsePreset(part string) (string, string, error) {
	eqParts := strings.SplitN(part, " = ", 2)
	if len(eqParts) != 2 {
		return "", "", fmt.Errorf("could not parse preset part %s", part)
	}

	objPath, preSet := eqParts[0], eqParts[1]

	return objPath, preSet, nil
}

func parseStrWithSeparators(s string, separators []string) ([]string, error) {
	parsedValues := []string{}
	forcedEnd := false
	remainingStr := s

	for _, separator := range separators {
		if strings.HasSuffix(separator, "$") {
			separator = separator[:len(separator)-1]
			forcedEnd = true
		}

		before, after, ok := strings.Cut(remainingStr, separator)
		if !ok {
			log.Printf("separator not found '%s' '%s', skipping", s, separator)
			break
		}

		if before != "" {
			parsedValues = append(parsedValues, before)
		}

		remainingStr = after

		if forcedEnd {
			remainingStr = ""
			break
		}
	}

	if remainingStr != "" {
		parsedValues = append(parsedValues, remainingStr)
	}

	return parsedValues, nil
}
