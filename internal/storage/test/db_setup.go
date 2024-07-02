package test

import (
	"github.com/emsi-zero/auth_ir/internal/conf"
	"github.com/emsi-zero/auth_ir/internal/storage"
)

func SetupDBConnection(globalConfig *conf.GlobalConfiguration) (*storage.Connection, error) {
	return storage.Dial(globalConfig)
}
