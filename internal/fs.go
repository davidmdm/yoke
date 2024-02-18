package internal

import (
	"os"

	"gopkg.in/yaml.v3"
)

func WriteYAML(filename string, value any) (err error) {
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer func() {
		if closeErr := file.Close(); closeErr != nil && err == nil {
			err = closeErr
		}
	}()

	encoder := yaml.NewEncoder(file)
	encoder.SetIndent(2)

	return encoder.Encode(value)
}
