package graph

import (
	"bluechat-server/database"
	"bluechat-server/graph/model"
)

// This file will not be regenerated automatically.
//
// It serves as dependency injection for your app, add any dependencies you require here.

type Resolver struct {
	Database database.Database
	ChatObservers map[string]chan []*model.ChatMessages
}
