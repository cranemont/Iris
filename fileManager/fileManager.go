package fileManager

import (
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
	os.Mkdir(f.basePath+"/"+name, os.FileMode(constants.BASE_FILE_MODE))
}

func (f *FileManager) RemoveDir(name string) {
	os.RemoveAll(f.basePath + "/" + name)
}
