package file

import (
	"bytes"
	"log"
	"os"

	"github.com/cranemont/judge-manager/constants"
)

func CreateDir(dir string) error {
	if err := os.Mkdir(constants.BASE_DIR+"/"+dir, os.FileMode(constants.BASE_FILE_MODE)); err != nil {
		log.Println(err)
		return err
	}
	return nil
}

func RemoveDir(dir string) error {
	if err := os.RemoveAll(constants.BASE_DIR + "/" + dir); err != nil {
		log.Println(err)
		return err
	}
	return nil
}

func CreateFile(path string, data string) error {
	if err := os.WriteFile(path, []byte(data), constants.BASE_FILE_MODE); err != nil {
		log.Println(err)
		return err
	}
	return nil
}

func ReadFile(path string) ([]byte, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	return data, nil
}

func MakeFilePath(dir string, fileName string) *bytes.Buffer {
	var b bytes.Buffer
	b.WriteString(constants.BASE_DIR)
	b.WriteString("/")
	b.WriteString(dir)
	b.WriteString("/")
	b.WriteString(fileName)
	return &b
}
