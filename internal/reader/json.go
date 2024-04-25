package reader

import (
	"encoding/json"
	"io"
	"os"
)

func ReadJsonFromFile(path string, target any) error {
	jsonFile, err := os.Open(path)
	if err != nil {
		return err
	}

	defer func(jsonFile *os.File) {
		_ = jsonFile.Close()
	}(jsonFile)

	byteValue, err := io.ReadAll(jsonFile)
	if err != nil {
		return err
	}

	err = json.Unmarshal(byteValue, target)
	if err != nil {
		return err
	}

	return nil
}
