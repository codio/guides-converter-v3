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

func removeByPath(path string) error {
	err := os.Remove(path)
  if err != nil {
    return err
  }
	return nil
}
