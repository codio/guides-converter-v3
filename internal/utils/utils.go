package utils

import (
	"encoding/json"
	"io/ioutil"
	"os"
)

func RemoveFile(file string) error {
	if err := os.Remove(file); err != nil && !os.IsNotExist(err) {
		return err
	}
	return nil
}

func RemoveDirectory(dir string) error {
	if err := os.RemoveAll(dir); err != nil && !os.IsNotExist(err) {
		return err
	}
	return nil
}

func RemoveDirectoryIfEmpty(dir string) error {
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		return nil
	}
	items, err := os.ReadDir(dir)
	if err != nil {
		return err
	}
	if len(items) == 0 {
		RemoveDirectory(dir)
	}
	return nil
}

func MakeDir(directory string) error {
	if stat, err := os.Stat(directory); err == nil && stat.IsDir() {
		return nil
	}
	if err := os.Mkdir(directory, 0777); err != nil {
		return err
	}
	return nil
}

func Rename(oldPath, newPath string) error {
	if err := os.Rename(oldPath, newPath); err != nil {
		return err
	}
	return nil
}

func GetParsedJson(path string, parsed any) error {
	jsonFile, err := os.OpenFile(path, os.O_RDONLY, 0)
	if err != nil {
		return err
	}
	defer jsonFile.Close()

	bytes, err := ioutil.ReadAll(jsonFile)
	if err != nil {
		return err
	}
	if err := json.Unmarshal(bytes, parsed); err != nil {
		return err
	}
	return nil
}

func WriteJson(path string, content map[string]interface{}) error {
	jsonFile, err := os.OpenFile(path, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		return err
	}
	defer jsonFile.Close()
	data, err := json.MarshalIndent(content, "", "  ")
	if err != nil {
		return err
	}
	if err := jsonFile.Truncate(0); err != nil {
		return err
	}
	if _, err := jsonFile.Seek(0, 0); err != nil {
		return err
	}
	if _, err := jsonFile.Write(data); err != nil {
		return err
	}
	return nil
}
