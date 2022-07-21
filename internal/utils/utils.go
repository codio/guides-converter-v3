package utils

import (
	"os"
)

func RemoveFile(file string) error {
  return removeByPath(file)
}

func RemoveDirectoryIfEmpty(dir string) error {
	items, err := os.ReadDir(dir)
	if err != nil {
		return err
	}
	if len(items) == 0 {
		removeByPath(dir)
	}
	return nil
}

func removeByPath(path string) error {
	err := os.Remove(path)
  if err != nil {
    return err
  }
	return nil
}
