package file

import (
	"github.com/cranemont/judge-manager/router"
	"github.com/cranemont/judge-manager/service/file"
	"github.com/cranemont/judge-manager/service/logger"
)

type connector struct {
	file   file.FileManager
	router router.Router
	logger logger.Logger
}

func NewConnector(f file.FileManager, r router.Router, l logger.Logger) *connector {
	return &connector{f, r, l}
}

func (c *connector) Connect() {
	// validate file source
}

func (c *connector) Disconnect() {
	// relaese resources
}

func (c *connector) Handle() {

}
