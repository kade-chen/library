package file

import (
	"encoding/json"
	"io"
	"os"

	"github.com/BurntSushi/toml"
	"sigs.k8s.io/yaml"
)

func ReadFile(path string) (string, error) {
	fs, err := os.Open(path)
	if err != nil {
		return "", err
	}
	content, err := io.ReadAll(fs)
	if err != nil {
		return "", err
	}

	return string(content), nil
}

func ReadContentFile(filepath string) ([]byte, error) {
	fd, err := os.Open(filepath)
	if err != nil {
		return nil, err
	}
	defer fd.Close()

	payload, err := io.ReadAll(fd)
	if err != nil {
		return nil, err
	}
	return payload, nil
}

func MustReadContentFile(filepath string) string {
	content, err := ReadContentFile(filepath)
	if err != nil {
		panic(err)
	}
	return string(content)
}

func ReadYamlFile(filepath string, v any) error {
	content, err := ReadContentFile(filepath)
	if err != nil {
		return err
	}
	return yaml.Unmarshal(content, v)
}

func MustReadYamlFile(filepath string, v any) {
	err := ReadYamlFile(filepath, v)
	if err != nil {
		panic(err)
	}
}

func MustReadJsonFile(filepath string, v any) {
	err := ReadJsonFile(filepath, v)
	if err != nil {
		panic(err)
	}
}

func ReadJsonFile(filepath string, v any) error {
	fd, err := os.Open(filepath)
	if err != nil {
		return err
	}
	defer fd.Close()

	payload, err := io.ReadAll(fd)
	if err != nil {
		return err
	}
	return json.Unmarshal(payload, v)
}

func MustToToml(key string, value any, path string) error {
	f, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0666)
	if err != nil {
		panic(err)
	}
	appConf := map[string]any{key: value}
	toml.NewEncoder(f).Encode(appConf)
	return nil
}
