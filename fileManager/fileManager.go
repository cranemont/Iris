package fileManager

import (
	"log"
	"os"

	"github.com/cranemont/judge-manager/constants"
)

type FileManager interface {
	CreateDir(name string) error
	RemoveDir(name string) error
	CreateFile(srcPath string, data string) error
}

type fileManager struct {
	basePath string
}

func NewFileManager() *fileManager {
	return &fileManager{basePath: constants.BASE_DIR}
}

func (f *fileManager) CreateDir(name string) error {
	if err := os.Mkdir(constants.BASE_DIR+"/"+name, os.FileMode(constants.BASE_FILE_MODE)); err != nil {
		return err
	}
	return nil
}

func (f *fileManager) RemoveDir(name string) error {
	if err := os.RemoveAll(constants.BASE_DIR + "/" + name); err != nil {
		log.Println(err)
		return err
	}
	return nil
}

func (f *fileManager) CreateFile(srcPath string, data string) error {
	if err := os.WriteFile(srcPath, []byte(data), constants.BASE_FILE_MODE); err != nil {
		return err
	}
	return nil
}
