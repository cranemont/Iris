package fileManager

import (
	"fmt"
	"os"

	"github.com/cranemont/judge-manager/constants"
)

type FileManager struct {
	basePath string
}

func NewFileManager() *FileManager {
	return &FileManager{basePath: constants.BASE_DIR}
}

func (f *FileManager) CreateDir(name string) {
	err := os.Mkdir(constants.BASE_DIR+"/"+name, os.FileMode(constants.BASE_FILE_MODE))
	if err != nil {
		fmt.Println(err)
	}
}

func (f *FileManager) RemoveDir(name string) {
	os.RemoveAll(constants.BASE_DIR + "/" + name)
}
