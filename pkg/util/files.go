package util

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"path/filepath"
)

func WriteJson(path string, data []byte) error {
	ensureDirectory(path)
	return ioutil.WriteFile(path, data, 0666)
}

func ensureDirectory(path string) error {
	parent := filepath.Dir(path)
	if _, err := os.Stat(parent); err != nil {
		return os.MkdirAll(parent, 0666)
	}
	return nil
}

func ReadFileAs(path string, obj interface{}) error {
	d, err := ioutil.ReadFile(path)
	if err != nil {
		return err
	}
	err = json.Unmarshal(d, obj)
	if err != nil {
		return err
	}
	return nil
}
