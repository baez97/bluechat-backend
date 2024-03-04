package graph

import (
	"bluechat-server/graph/model"
	"database/sql"
)

// This file will not be regenerated automatically.
//
// It serves as dependency injection for your app, add any dependencies you require here.

type Resolver struct {
	Database *sql.DB
	ChatObservers map[string]chan []*model.Message
}
