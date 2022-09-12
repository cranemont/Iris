package file

import (
	"bytes"
	"log"
	"os"

	"github.com/cranemont/judge-manager/constants"
)

type FileManager interface {
	CreateDir(dir string) error
	RemoveDir(dir string) error
	CreateFile(path string, data string) error
	ReadFile(path string) ([]byte, error)
	MakeFilePath(dir string, fileName string) *bytes.Buffer
}

type fileManager struct {
	baseDir string
}

func NewFileManager() *fileManager {
	fileManager := fileManager{}
	if os.Getenv("APP_ENV") == "production" {
		fileManager.baseDir = constants.OUTPUT_PATH_PROD
	} else {
		fileManager.baseDir = constants.OUTPUT_PATH_DEV
	}
	return &fileManager
}

func (f *fileManager) CreateDir(dir string) error {
	if err := os.Mkdir(f.baseDir+"/"+dir, os.FileMode(constants.BASE_FILE_MODE)); err != nil {
		log.Println(err)
		return err
	}
	return nil
}

func (f *fileManager) RemoveDir(dir string) error {
	if err := os.RemoveAll(f.baseDir + "/" + dir); err != nil {
		log.Println(err)
		return err
	}
	return nil
}

func (f *fileManager) CreateFile(path string, data string) error {
	if err := os.WriteFile(path, []byte(data), constants.BASE_FILE_MODE); err != nil {
		log.Println(err)
		return err
	}
	return nil
}

func (f *fileManager) ReadFile(path string) ([]byte, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	return data, nil
}

func (f *fileManager) MakeFilePath(dir string, fileName string) *bytes.Buffer {
	var b bytes.Buffer
	b.WriteString(f.baseDir)
	b.WriteString("/")
	b.WriteString(dir)
	b.WriteString("/")
	b.WriteString(fileName)
	return &b
}
