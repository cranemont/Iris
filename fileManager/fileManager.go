package fileManager

import (
	"fmt"
	"os"

	"github.com/cranemont/judge-manager/constants"
)

type FileManager interface {
	CreateDir(name string)
	RemoveDir(name string)
}

type fileManager struct {
	basePath string
}

func NewFileManager() *fileManager {
	return &fileManager{basePath: constants.BASE_DIR}
}

func (f *fileManager) CreateDir(name string) {
	err := os.Mkdir(constants.BASE_DIR+"/"+name, os.FileMode(constants.BASE_FILE_MODE))
	if err != nil {
		fmt.Println(err)
	}
}

func (f *fileManager) RemoveDir(name string) {
	os.RemoveAll(constants.BASE_DIR + "/" + name)
}
