package utils

import (
	"fmt"
	"os"
)

func PathExists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}

func CreateDir(dirs ...string) (err error) {
	for _, v := range dirs {
		exist, err := PathExists(v)
		if err != nil {
			return err
		}
		if !exist {
			fmt.Print("create directory" + v)
			if err := os.MkdirAll(v, os.ModePerm); err != nil {
				fmt.Print("create directory" + v)
				return err
			}
		}
	}
	return err
}
