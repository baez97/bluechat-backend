package database

import (
	"bluechat-server/graph/model"
	"database/sql"
	"fmt"
)

type Database struct {
	SQL *sql.DB
}

func (db *Database) GetMessagesSince(userId string, since string) ([]*model.Message, error){
	query := `
		SELECT id, sender_id, receiver_id, content, media_url, timestamp
		FROM Messages
		WHERE receiver_id = $1 AND timestamp > $2
		ORDER BY timestamp DESC
	`
	rows, err := db.SQL.Query(query, userId, since)
	if err != nil {
		fmt.Printf("error querying database: %v\n", err)
		return nil, err
	}
	defer rows.Close()
	messages := make([]*model.Message, 0)
	for rows.Next() {
		// Read the message data from the database
		var message model.Message
		err := rows.Scan(&message.ID, &message.SenderID, &message.ReceiverID, &message.Content, &message.MediaURL, &message.Timestamp)
		if err != nil {
			fmt.Printf("error scanning row: %v\n", err)
			return nil, err
		}
		
		messages = append(messages, &message)
	}
	return messages, nil
}