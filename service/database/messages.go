package database

import (
	"database/sql"
	"fmt"
)

func (db *appdbimpl) GetMessage(id_message int64) (Message, error) {
	var message Message

	// COALESCE fixa:
	// error getting message: database error getting message: sql: Scan error on column index 5, name "forwarded_source": converting NULL to int64 is unsupported

	err := db.c.QueryRow(`SELECT id, text, photo, author, recipient, COALESCE(forwarded_source, 0) AS forwarded_source, timestamp FROM messages WHERE id = $1;`, id_message).Scan(
		&message.Id, &message.Content, &message.Photo, &message.Author, &message.Recipient, &message.Forwarded_source, &message.Timestamp)
	if err == sql.ErrNoRows {
		return Message{}, fmt.Errorf("message not found")
	}
	if err != nil {
		return Message{}, fmt.Errorf("database error getting message: %w", err)
	}

	return message, nil
}

func (db *appdbimpl) GetLastMessage(id_conversation int64) (Message, error) {
	var message Message

	err := db.c.QueryRow(`
		SELECT
			id,
			text,
			photo,
			author,
			recipient,
			COALESCE(forwarded_source, 0) AS forwarded_source,
			timestamp
		FROM
			messages
		WHERE
			recipient = $1
		ORDER BY
			timestamp DESC
		LIMIT 1
		`, id_conversation).Scan(
		&message.Id,
		&message.Content,
		&message.Photo,
		&message.Author,
		&message.Recipient,
		&message.Forwarded_source,
		&message.Timestamp,
	)
	if err == sql.ErrNoRows {
		return Message{}, nil
	}
	if err != nil {
		return Message{}, fmt.Errorf("error getting last message: %w", err)
	}

	return message, nil
}

// ############################################################

func (db *appdbimpl) GetMessagesOfConversation(id_conversation int64, id_auth int64) ([]Message, error) {
	// Check if the conversation exists
	_, err := db.GetConversation(id_conversation)
	if err != nil && err.Error() == "no conversation found" {
		return []Message{}, err
	}
	if err != nil {
		return []Message{}, fmt.Errorf("error getting conversation: %w", err)
	}

	// Check if the user is in the conversation
	isUserInConversation, err := db.IsUserInConversation(id_conversation, id_auth)
	if err != nil {
		return []Message{}, fmt.Errorf("database error checking if user is in conversation: %w", err)
	}
	if !isUserInConversation {
		return []Message{}, fmt.Errorf("user is not in conversation")
	}

	rows, err := db.c.Query(`
		SELECT
			id,
			text,
			photo,
			author,
			recipient,
			COALESCE(forwarded_source, 0) AS forwarded_source,
			timestamp
		FROM
			messages
		WHERE
			recipient = $1
		ORDER BY
			timestamp ASC
		`, id_conversation)
	if err != nil {
		return []Message{}, fmt.Errorf("database error getting messages of conversation: %w", err)
	}
	defer rows.Close()

	var messages []Message
	for rows.Next() {
		var message Message
		err := rows.Scan(
			&message.Id,
			&message.Content,
			&message.Photo,
			&message.Author,
			&message.Recipient,
			&message.Forwarded_source,
			&message.Timestamp,
		)
		if err != nil {
			return []Message{}, fmt.Errorf("database error scanning messages of conversation: %w", err)
		}
		messages = append(messages, message)
	}

	return messages, nil
}

func (db *appdbimpl) SendMessage(id_conversation int64, id_auth int64, message Message) (Message, error) {

	// Check if the conversation exists
	_, err := db.GetConversation(id_conversation)
	if err != nil && err.Error() == "no conversation found" {
		return Message{}, err
	}
	if err != nil {
		return Message{}, fmt.Errorf("error getting conversation: %w", err)
	}

	// Check if the user is in the conversation
	isUserInConversation, err := db.IsUserInConversation(id_conversation, id_auth)
	if err != nil {
		return Message{}, fmt.Errorf("database error checking if user is in conversation: %w", err)
	}
	if !isUserInConversation {
		return Message{}, fmt.Errorf("user is not in conversation")
	}

	message.Author = id_auth
	message.Recipient = id_conversation

	res, err := db.c.Exec(`
	INSERT INTO messages (
	text, 
	photo, 
	author, 
	recipient
	) VALUES ($1, $2, $3, $4);`, message.Content, message.Photo, message.Author, message.Recipient)
	if err != nil {
		return Message{}, fmt.Errorf("database error inserting message: %w", err)
	}

	message.Id, err = res.LastInsertId()
	if err != nil {
		return Message{}, fmt.Errorf("database error getting last inserted id of added message: %w", err)
	}

	message, err = db.GetMessage(message.Id)
	if err != nil {
		return Message{}, fmt.Errorf("database error getting message: %w", err)
	}

	return message, nil
}

func (db *appdbimpl) DeleteMessage(id_message int64, id_auth int64) error {
	// Check if the message existss
	message, err := db.GetMessage(id_message)
	if err != nil && err.Error() == "message not found" {
		return err
	}
	if err != nil {
		return fmt.Errorf("error getting message: %w", err)
	}

	// Check if the user is the author of the message
	if message.Author != id_auth {
		return fmt.Errorf("user is not the author of the message")
	}

	_, err = db.c.Exec(`DELETE FROM messages WHERE id = $1;`, id_message)
	if err != nil {
		return fmt.Errorf("database error deleting message: %w", err)
	}

	return nil
}

func (db *appdbimpl) ForwardMessage(id_message int64, id_auth int64, id_conversation int64) (Message, error) {

	// Check if the message existss
	message, err := db.GetMessage(id_message)
	if err != nil && err.Error() == "message not found" {
		return Message{}, err
	}
	if err != nil {
		return Message{}, fmt.Errorf("error getting message: %w", err)
	}

	// // Check if the user is the author of the message
	// if message.Author != id_auth {
	// 	return Message{}, fmt.Errorf("user is not the author of the message")
	// }

	// Check if the user is in the forwoarded conversation
	isUserInConversation, err := db.IsUserInConversation(id_conversation, id_auth)
	if err != nil {
		return Message{}, fmt.Errorf("database error checking if user is in conversation: %w", err)
	}
	if !isUserInConversation {
		return Message{}, fmt.Errorf("user is not in conversation")
	}

	message.Author = id_auth
	message.Forwarded_source = id_message

	res, err := db.c.Exec(`
	INSERT INTO messages (
	text,
	photo, 
	author, 
	recipient,
	forwarded_source
	) VALUES ($1, $2, $3, $4, $5);`, message.Content, message.Photo, message.Author, message.Recipient, message.Forwarded_source)
	if err != nil {
		return Message{}, fmt.Errorf("database error inserting message: %w", err)
	}

	message.Id, err = res.LastInsertId()
	if err != nil {
		return Message{}, fmt.Errorf("database error getting last inserted id of added message: %w", err)
	}

	message, err = db.GetMessage(message.Id)
	if err != nil {
		return Message{}, fmt.Errorf("database error getting message: %w", err)
	}

	return message, nil
}
