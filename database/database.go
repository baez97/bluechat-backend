package database

import (
	"bluechat-server/graph/model"
	"database/sql"
	"encoding/json"
	"fmt"
)

type Database struct {
	SQL *sql.DB
}

func (db *Database) SetUserTimeStamp(userId string, timestamp string) error {
	query := `
		UPDATE Users
		SET timestamp=$1
		WHERE id=$2
	`
	stmt, err := db.SQL.Prepare(query)
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(timestamp, userId)
	if err != nil {
		return err
	}

	return nil
}
func (db *Database) GetUserTimestamp(userId string) (string, error){
	var timestamp string
	query := "SELECT timestamp FROM users WHERE id = $1"
	err := db.SQL.QueryRow(query, userId).Scan(&timestamp)
	if err != nil {
			// If there's an error, return an empty string and the error
			return "", err
	}
	// If successful, return the timestamp and nil error
	return timestamp, nil
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

func (db *Database) GetChatMessagesSince(userId string, since string) ([]*model.ChatMessages, error){
	query := `
		SELECT
			sender_id,
			json_agg(
					json_build_object(
							'id', id::TEXT,
							'senderId', sender_id::TEXT,
							'receiverId', receiver_id::TEXT,
							'content', content,
							'media_url', media_url,
							'timestamp', timestamp
					) ORDER BY timestamp DESC
			) AS messages
		FROM
			(
					SELECT
							sender_id,
							receiver_id,
							id,
							content,
							media_url,
							timestamp
					FROM
							Messages
					WHERE
							receiver_id = $1
							AND timestamp > $2

					UNION ALL

					SELECT
							receiver_id AS sender_id,
							receiver_id AS receiver_id,
							id,
							content,
							media_url,
							timestamp
					FROM
							Messages
					WHERE
							sender_id = $1
							AND timestamp > $2
			) AS combined_messages
		GROUP BY
			sender_id
		ORDER BY
			MAX(timestamp) DESC;
	`
	rows, err := db.SQL.Query(query, userId, since)
	if err != nil {
		fmt.Printf("error querying database: %v\n", err)
		return nil, err
	}
	defer rows.Close()
	chatMessages := make([]*model.ChatMessages, 0)
	for rows.Next() {
		// Read the message data from the database
		var senderId string
		var messagesJSON string
		err := rows.Scan(&senderId, &messagesJSON)
		if err != nil {
			fmt.Printf("error scanning row: %v\n", err)
			return nil, err
		}
		var messages []*model.Message
		fmt.Println(messagesJSON)
		err = json.Unmarshal([]byte(messagesJSON), &messages)
		if err != nil {
			fmt.Printf("error unmarshalling messages: %v\n", err)
			return nil, err
		}
		chatMessages = append(chatMessages, &model.ChatMessages{
			SenderID: senderId,
			Messages: messages,
		})
	}
	return chatMessages, nil
}

func (db *Database) GetUsers(_ string) ([]*model.User, error) {
	query := "SELECT id, username, display_name, photo_url, company_id, timestamp FROM users"
	rows, err := db.SQL.Query(query)
	if err != nil {
		fmt.Printf("error querying database: %v\n", err)
		return nil, err
	}
	defer rows.Close()
	users := []*model.User{}
	for rows.Next() {
		var user model.User
		err = rows.Scan(&user.ID, &user.Username, &user.DisplayName, &user.PhotoURL, &user.CompanyID, &user.Timestamp)
		if err != nil {
			fmt.Printf("error scanning row: %v\n", err)
			return nil, err
		}
		users = append(users, &user)
	}
	return users, nil
}