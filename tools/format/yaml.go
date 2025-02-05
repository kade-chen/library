package format

import (
	"encoding/json"

	"sigs.k8s.io/yaml"
)


func MustToYaml(v any) string {
	b, err := yaml.Marshal(v)
	if err != nil {
		panic(err)
	}
	return string(b)
}

// This function takes in a byte slice and an empty interface and returns an error
func Unmarshal(data []byte, v any) error {
	// Unmarshal the byte slice into the empty interface
	return json.Unmarshal(data, v)
}