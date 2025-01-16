package common

import "os"

func CreateDir(dirPath string) error {
	if err := os.MkdirAll(dirPath, 0755); err != nil && !os.IsExist(err) {
		return err
	}
	return nil
}
