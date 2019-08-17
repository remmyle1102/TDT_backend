package utils

import (
	"io/ioutil"
	"log"
	"os"

	"github.com/sirupsen/logrus"
)

// ListDirOrFile listType 1 means listing dir, 2 means listing files
func ListDirOrFile(dir string, listType int) ([]string, error) {
	dirContents, err := ioutil.ReadDir(dir)
	if err != nil {
		log.Fatal(err)
	}
	result := make([]string, 0)
	switch listType {
	case 1:
		for _, f := range dirContents {
			if f.IsDir() {
				result = append(result, f.Name())
			}
		}

	case 2:
		for _, f := range dirContents {
			if f.IsDir() == false {
				result = append(result, f.Name())
			}
		}
	}
	return result, err
}

// CreateDirIfNotExist create directory if not exist
func CreateDirIfNotExist(dir string) error {
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		err = os.MkdirAll(dir, 0755)
		if err != nil {
			logrus.Error(err)
			return err
		}
	}
	return nil
}
