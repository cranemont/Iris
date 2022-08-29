package fileManager

import (
	"io/ioutil"
	"os"

	"github.com/cranemont/judge-manager/constants"
)

type FileManager interface {
	CreateDir(name string) error
	RemoveDir(name string)
	CreateFile(srcPath string, data string) error
}

type fileManager struct {
	basePath string
}

func NewFileManager() *fileManager {
	return &fileManager{basePath: constants.BASE_DIR}
}

func (f *fileManager) CreateDir(name string) error {
	err := os.Mkdir(constants.BASE_DIR+"/"+name, os.FileMode(constants.BASE_FILE_MODE))
	if err != nil {
		return err
	}
	return nil
}

func (f *fileManager) RemoveDir(name string) {
	os.RemoveAll(constants.BASE_DIR + "/" + name)
}

func (f *fileManager) CreateFile(srcPath string, data string) error {
	err := ioutil.WriteFile(srcPath, []byte(data), constants.BASE_FILE_MODE)
	if err != nil {
		return err
	}
	return nil
}
