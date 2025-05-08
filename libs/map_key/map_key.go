package mapkey

import (
	"fmt"
	"strconv"
	"strings"
)

type Result struct {
	Value    any
	FullPath string
}

// Examples of valid paths:
// - "users.*.name"         // get all user names
// - "company.*.manager.name"  // get all manager names across departments
// - "data.items.*.tags.*"  // get all tags from all items
// - "users.0.name"         // get the first user name
func Get(data any, path string) ([]Result, error) {
	if path == "" {
		return []Result{{Value: data, FullPath: ""}}, nil
	}

	parts := strings.Split(path, ".")
	return getRecursive(data, parts, "", []Result{})
}

func getRecursive(data any, parts []string, currentPath string, results []Result) ([]Result, error) {
	if len(parts) == 0 {
		fullPath := strings.TrimPrefix(currentPath, ".")
		results = append(results, Result{Value: data, FullPath: fullPath})
		return results, nil
	}

	part := parts[0]
	remainingParts := parts[1:]

	switch current := data.(type) {
	case map[string]any:
		if part == "*" {
			for key, value := range current {
				newPath := currentPath
				if newPath != "" {
					newPath += "."
				}
				newPath += key
				var err error
				results, err = getRecursive(value, remainingParts, newPath, results)
				if err != nil {
					return nil, err
				}
			}
			return results, nil
		}

		value, exists := current[part]
		if !exists {
			return nil, fmt.Errorf("key not found: %s", part)
		}
		newPath := currentPath
		if newPath != "" {
			newPath += "."
		}
		newPath += part
		return getRecursive(value, remainingParts, newPath, results)

	case []any:
		if part == "*" {
			for i, value := range current {
				newPath := currentPath
				if newPath != "" {
					newPath += "."
				}
				newPath += strconv.Itoa(i)
				var err error
				results, err = getRecursive(value, remainingParts, newPath, results)
				if err != nil {
					return nil, err
				}
			}
			return results, nil
		}

		index, err := strconv.Atoi(part)
		if err != nil {
			return nil, fmt.Errorf("invalid array index: %s", part)
		}
		if index < 0 || index >= len(current) {
			return nil, fmt.Errorf("array index out of bounds: %d", index)
		}
		newPath := currentPath
		if newPath != "" {
			newPath += "."
		}
		newPath += strconv.Itoa(index)
		return getRecursive(current[index], remainingParts, newPath, results)

	default:
		return nil, fmt.Errorf("cannot navigate through type %T at path segment: %s", current, part)
	}
}

func GetStrings(data any, path string) (map[string]string, error) {
	results, err := Get(data, path)
	if err != nil {
		return nil, err
	}

	stringResults := make(map[string]string)
	for _, result := range results {
		str, ok := result.Value.(string)
		if !ok {
			return nil, fmt.Errorf("value at path %s is not a string", result.FullPath)
		}
		stringResults[result.FullPath] = str
	}
	return stringResults, nil
}

func GetInts(data any, path string) (map[string]int, error) {
	results, err := Get(data, path)
	if err != nil {
		return nil, err
	}

	intResults := make(map[string]int)
	for _, result := range results {
		num, ok := result.Value.(int)
		if !ok {
			return nil, fmt.Errorf("value at path %s is not an integer", result.FullPath)
		}
		intResults[result.FullPath] = num
	}
	return intResults, nil
}

func GetFloats(data any, path string) (map[string]float64, error) {
	results, err := Get(data, path)
	if err != nil {
		return nil, err
	}

	floatResults := make(map[string]float64)
	for _, result := range results {
		num, ok := result.Value.(float64)
		if !ok {
			return nil, fmt.Errorf("value at path %s is not a float", result.FullPath)
		}
		floatResults[result.FullPath] = num
	}
	return floatResults, nil
}
