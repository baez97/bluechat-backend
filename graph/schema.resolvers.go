package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.
// Code generated by github.com/99designs/gqlgen version v0.17.44

import (
	"bluechat-server/graph/model"
	"context"
	"fmt"
	"time"
)

// PostMessage is the resolver for the postMessage field.
func (r *mutationResolver) PostMessage(ctx context.Context, senderID string, receiverID *string, groupID *string, content string, mediaURL *string) (*model.Message, error) {
	// Determine if the message is for a user or a group
	if receiverID == nil && groupID == nil {
		return nil, fmt.Errorf("receiverID or groupID must be provided")
	}

	// Insert the new message into the database
	query := `
		INSERT INTO Messages (sender_id, receiver_id, group_id, content, media_url, timestamp)
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING id
	`
	var messageID string
	err := r.Database.QueryRowContext(ctx, query, senderID, receiverID, groupID, content, mediaURL, time.Now()).Scan(&messageID)
	if err != nil {
		return nil, fmt.Errorf("error inserting message into database: %v", err)
	}

	// Create a message object with the inserted data
	message := &model.Message{
		ID:         messageID,
		SenderID:   senderID,
		Content:    &content,
		MediaURL:   mediaURL,
		Timestamp:  time.Now().Format(time.RFC3339), // Convert timestamp to string
	}

	if receiverID != nil {
		message.ReceiverID = *receiverID
	}
	if groupID != nil {
		message.GroupID = *groupID
	}

	messages := []*model.Message{message}
	if (r.ChatObservers[*receiverID] != nil) {
		r.ChatObservers[*receiverID] <- messages
	}

	return message, nil
}

// CreateGroup is the resolver for the createGroup field.
func (r *mutationResolver) CreateGroup(ctx context.Context, name string, userIds []string) (*model.Group, error) {
	panic(fmt.Errorf("not implemented: CreateGroup - createGroup"))
}

// AddUserToGroup is the resolver for the addUserToGroup field.
func (r *mutationResolver) AddUserToGroup(ctx context.Context, userID string, groupID string) (*model.Group, error) {
	panic(fmt.Errorf("not implemented: AddUserToGroup - addUserToGroup"))
}

// DeleteUserFromGroup is the resolver for the deleteUserFromGroup field.
func (r *mutationResolver) DeleteUserFromGroup(ctx context.Context, userID string, groupID string) (*model.Group, error) {
	panic(fmt.Errorf("not implemented: DeleteUserFromGroup - deleteUserFromGroup"))
}

// ModifyGroupUsers is the resolver for the modifyGroupUsers field.
func (r *mutationResolver) ModifyGroupUsers(ctx context.Context, groupID string, userIds []string) (*model.Group, error) {
	panic(fmt.Errorf("not implemented: ModifyGroupUsers - modifyGroupUsers"))
}

// CreateCompany is the resolver for the createCompany field.
func (r *mutationResolver) CreateCompany(ctx context.Context, name string, photoURL *string) (*model.Company, error) {
	panic(fmt.Errorf("not implemented: CreateCompany - createCompany"))
}

// Group is the resolver for the group field.
func (r *queryResolver) Group(ctx context.Context, id string) (*model.Group, error) {
	panic(fmt.Errorf("not implemented: Group - group"))
}

// Messages is the resolver for the Messages field.
func (r *queryResolver) Messages(ctx context.Context, userID string, since string) ([]*model.Message, error) {
	query := `
		SELECT id, sender_id, receiver_id, content, media_url, timestamp
		FROM Messages
		WHERE receiver_id=$1 AND timestamp > $2
		ORDER BY timestamp DESC	
	`
	rows, err := r.Database.Query(query, userID, since)
	if err != nil {
		fmt.Printf("error querying database: %v\n", err)
		return nil, err
	}
	defer rows.Close()

	messages := []*model.Message{}
	for rows.Next() {
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

// NewMessages is the resolver for the NewMessages field.
func (r *subscriptionResolver) NewMessages(ctx context.Context, userID string, since string) (<-chan []*model.Message, error) {
	// Create a channel to send messages to the client
	messageChan := make(chan []*model.Message)

	// Start a goroutine to listen for new messages and send them to the client
	go func() {
		// Set up a SQL query to select new messages for the specified user
		query := `
			SELECT id, sender_id, receiver_id, content, media_url, timestamp
			FROM Messages
			WHERE receiver_id = $1 AND timestamp > $2
			ORDER BY timestamp DESC
		`

		rows, err := r.Database.Query(query, userID, since)
		if err != nil {
			fmt.Printf("error querying database: %v\n", err)
			return
		}
		defer rows.Close()

		messages := make([]*model.Message, 0)
		for rows.Next() {
			// Read the message data from the database
			var message model.Message
			err := rows.Scan(&message.ID, &message.SenderID, &message.ReceiverID, &message.Content, &message.MediaURL, &message.Timestamp)
			if err != nil {
				fmt.Printf("error scanning row: %v\n", err)
				return
			}

			messages = append(messages, &message)
		}

		messageChan <- messages

		if err := rows.Err(); err != nil {
			fmt.Printf("error iterating over rows: %v\n", err)
			return
		}
	}()

	r.ChatObservers[userID] = messageChan

	return messageChan, nil
}

// Mutation returns MutationResolver implementation.
func (r *Resolver) Mutation() MutationResolver { return &mutationResolver{r} }

// Query returns QueryResolver implementation.
func (r *Resolver) Query() QueryResolver { return &queryResolver{r} }

// Subscription returns SubscriptionResolver implementation.
func (r *Resolver) Subscription() SubscriptionResolver { return &subscriptionResolver{r} }

type mutationResolver struct{ *Resolver }
type queryResolver struct{ *Resolver }
type subscriptionResolver struct{ *Resolver }