package database

import (
	"database/sql"
	"errors"
	"fmt"
)

func (db *appdbimpl) SetMessageReadByUser(id_message int64, id_user int64) error {

	isRead, err := db.IsMessageReadByUser(id_message, id_user)
	if err != nil {
		return fmt.Errorf("database error checking if message is read by user: %w", err)
	}

	if isRead {
		return nil
	}

	_, err = db.c.Exec(`
	INSERT INTO messages_readers (
		id_message,
		id_user
	) VALUES ($1, $2);`, id_message, id_user)
	if err != nil {
		return fmt.Errorf("database error setting message read: %w", err)
	}

	return nil
}

func (db *appdbimpl) IsMessageReadByUser(id_message int64, id_user int64) (bool, error) {
	var count int64
	err := db.c.QueryRow(`
	SELECT
		COUNT(*)
	FROM
		messages_readers
	WHERE
		id_message = $1
	AND
		id_user = $2;
	`, id_message, id_user).Scan(&count)
	if err != nil {
		return false, err
	}

	if count == 0 {
		return false, nil
	} else {
		return true, nil
	}
}

func (db *appdbimpl) GetMessageReadness(id_nessage int64, id_conversation int64) (bool, error) {
	var conversation Conversation

	conversation, err := db.GetConversation(id_conversation)
	if err != nil {
		return false, fmt.Errorf("database error getting conversation: %w", err)
	}

	var isRead bool = true

	for i := 0; i < len(conversation.Participants); i++ {
		// println(isRead)
		if !isRead {
			continue
		}

		isRead, err = db.IsMessageReadByUser(id_nessage, conversation.Participants[i].Id)
		if err != nil {
			return false, fmt.Errorf("database error checking if message is read by user: %w", err)
		}
	}

	return isRead, nil

}

func (db *appdbimpl) IsMessageReadInConversation(message Message) (bool, error) {

	var id_conversation int64 = message.Recipient

	conversation, err := db.GetConversation(id_conversation)
	if err != nil {
		return false, fmt.Errorf("database error getting conversation: %w", err)
	}

	for i := 0; i < len(conversation.Participants); i++ {
		isRead, err := db.IsMessageReadByUser(message.Id, conversation.Participants[i].Id)
		if err != nil {
			return false, fmt.Errorf("database error checking if message is read by user: %w", err)
		}
		if !isRead {
			return false, nil
		}
	}

	return true, nil
}

func (db *appdbimpl) GetReactionOfMessage(id_message int64) ([]Reaction, error) {
	rows, err := db.c.Query(`
		SELECT
			id_user,
			reaction
		FROM
			reactions
		WHERE
			id_message = $1;
	`, id_message)
	if err != nil {
		return []Reaction{}, fmt.Errorf("database error getting reactions: %w", err)
	}
	defer rows.Close()

	var reactions []Reaction
	for rows.Next() {
		var reaction Reaction
		err := rows.Scan(
			&reaction.User,
			&reaction.Reaction,
		)
		if err != nil {
			return []Reaction{}, fmt.Errorf("database error scanning reactions: %w", err)
		}
		reactions = append(reactions, reaction)
	}
	err = rows.Err()
	if err != nil {
		return []Reaction{}, fmt.Errorf("error getting conversation row: %w", err)
	}

	return reactions, nil

}

func (db *appdbimpl) GetMessage(id_message int64) (Message, error) {
	var message Message

	// COALESCE fixa:
	// error getting message: database error getting message: sql: Scan error on column index 5, name "forward": converting NULL to int64 is unsupported

	err := db.c.QueryRow(`
	SELECT 
		id, 
		text, 
		photo, 
		author, 
		recipient, 
		COALESCE(forward, 0) AS forward, 
		COALESCE(reply, 0) AS reply, 
		timestamp 
	FROM 
		messages 
	WHERE 
		id = $1;
	`, id_message).Scan(
		&message.Id,
		&message.Text,
		&message.Photo,
		&message.Author,
		&message.Recipient,
		&message.Forward,
		&message.Reply,
		&message.Timestamp,
	)

	if errors.Is(err, sql.ErrNoRows) {
		return Message{}, fmt.Errorf(ErrMessageNotFound)
	}
	if err != nil {
		return Message{}, fmt.Errorf("database error getting message: %w", err)
	}

	message.Reactions, err = db.GetReactionOfMessage(message.Id)
	if err != nil {
		return Message{}, fmt.Errorf("database error getting reactions: %w", err)
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
			COALESCE(forward, 0) AS forward,
			COALESCE(reply, 0) AS reply,
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
		&message.Text,
		&message.Photo,
		&message.Author,
		&message.Recipient,
		&message.Forward,
		&message.Reply,
		&message.Timestamp,
	)
	if errors.Is(err, sql.ErrNoRows) {
		return Message{}, nil
	}
	if err != nil {
		return Message{}, fmt.Errorf("error getting last message: %w", err)
	}

	message.Reactions, err = db.GetReactionOfMessage(message.Id)
	if err != nil {
		return Message{}, fmt.Errorf("database error getting reactions: %w", err)
	}

	return message, nil
}

func (db *appdbimpl) IsMessageInConversation(id_conversation int64, id_message int64) (bool, error) {
	var count int64

	err := db.c.QueryRow(`
	SELECT
	COUNT(*)
	FROM
	messages
	WHERE
	recipient = $1
	AND
	id = $2;
	`, id_conversation, id_message).Scan(&count)
	if err != nil {
		return false, err
	}

	if count == 0 {
		return false, nil
	} else {
		return true, nil
	}

}

func (db *appdbimpl) HasReaction(id_message int64, id_user int64) (bool, error) {
	var count int64
	err := db.c.QueryRow(`
		SELECT
			COUNT(*)
		FROM
			reactions
		WHERE
			id_user = $1
		AND
			id_message = $2;
		`, id_user, id_message).Scan(&count)
	if err != nil {
		return false, err
	}
	if count != 0 {
		return true, nil
	} else {
		return false, nil
	}
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
		return []Message{}, fmt.Errorf(ErrUserNotInConversation)
	}

	rows, err := db.c.Query(`
		SELECT
			id,
			text,
			photo,
			author,
			recipient,
			COALESCE(forward, 0) AS forward,
			COALESCE(reply, 0) AS reply,
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
			&message.Text,
			&message.Photo,
			&message.Author,
			&message.Recipient,
			&message.Forward,
			&message.Reply,
			&message.Timestamp,
		)
		if err != nil {
			return []Message{}, fmt.Errorf("database error scanning messages of conversation: %w", err)
		}
		messages = append(messages, message)
	}
	err = rows.Err()
	if err != nil {
		return []Message{}, fmt.Errorf("error getting conversation row: %w", err)
	}

	for i := 0; i < len(messages); i++ {
		messages[i].Reactions, err = db.GetReactionOfMessage(messages[i].Id)
		if err != nil {
			return []Message{}, fmt.Errorf("database error getting reactions: %w", err)
		}

		err = db.SetMessageReadByUser(messages[i].Id, id_auth)
		if err != nil {
			return []Message{}, fmt.Errorf("database error setting message read: %w", err)
		}

		messages[i].Read, err = db.GetMessageReadness(messages[i].Id, id_conversation)
		if err != nil {
			return []Message{}, fmt.Errorf("database error checking if message is read by user: %w", err)
		}
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
		return Message{}, fmt.Errorf(ErrUserNotInConversation)
	}

	// Check if reply message is in the conversation
	isMessageInConversation, err := db.IsMessageInConversation(id_conversation, message.Reply)
	if err != nil {
		return Message{}, err
	}
	if !isMessageInConversation && message.Reply != 0 {
		return Message{}, fmt.Errorf("message is not in the conversation")
	}

	message.Author = id_auth
	message.Recipient = id_conversation

	res, err := db.c.Exec(`
	INSERT INTO messages (
	text, 
	photo, 
	author, 
	recipient,
	reply
	) VALUES ($1, $2, $3, $4, $5);`, message.Text, message.Photo, message.Author, message.Recipient, message.Reply)
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
	if err != nil && err.Error() == ErrMessageNotFound {
		return err
	}
	if err != nil {
		return fmt.Errorf("error getting message: %w", err)
	}

	// Check if the user is the author of the message
	if message.Author != id_auth {
		return fmt.Errorf("user is not the author of the message")
	}

	// Check if the message has replies
	var replied_messages []int64

	rows, err := db.c.Query(`SELECT id FROM messages WHERE reply = $1;`, id_message)
	if err != nil {
		return fmt.Errorf("database error getting replies of message: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var id int64
		err := rows.Scan(&id)
		if err != nil {
			return fmt.Errorf("database error scanning replies of message: %w", err)
		}
		replied_messages = append(replied_messages, id)
	}
	err = rows.Err()
	if err != nil {
		return fmt.Errorf("error getting message row: %w", err)
	}

	_, _ = db.c.Exec(`SET FOREIGN_KEY_CHECKS = 0;`)

	// Remove replies
	for i := 0; i < len(replied_messages); i++ {
		_, err = db.c.Exec(`UPDATE messages SET reply = 0 WHERE id = $1;`, replied_messages[i])
		if err != nil {
			return fmt.Errorf("database error editing message: %w", err)
		}
	}

	_, _ = db.c.Exec(`SET FOREIGN_KEY_CHECKS = 1;`)

	// Check if the message has reactions
	hasReactions, err := db.HasReaction(id_message, id_auth)
	if err != nil {
		return fmt.Errorf("database error checking if message has reactions: %w", err)
	}
	if hasReactions {
		_, err = db.c.Exec(`DELETE FROM reactions WHERE id_message = $1;`, id_message)
		if err != nil {
			return fmt.Errorf("database error deleting reactions of message: %w", err)
		}
	}

	// Delete message
	_, err = db.c.Exec(`DELETE FROM messages WHERE id = $1;`, id_message)
	if err != nil {
		return fmt.Errorf("database error deleting message: %w", err)
	}

	return nil
}

func (db *appdbimpl) ForwardMessage(id_message int64, id_auth int64, id_conversation int64) (Message, error) {

	// Check if the message existss
	message, err := db.GetMessage(id_message)
	if err != nil && err.Error() == ErrMessageNotFound {
		return Message{}, err
	}
	if err != nil {
		return Message{}, fmt.Errorf("error getting message: %w", err)
	}

	// Check if the user is in the source conversation
	isUserInConversation, err := db.IsUserInConversation(message.Recipient, id_auth)
	if err != nil {
		return Message{}, fmt.Errorf("database error checking if user is in conversation: %w", err)
	}
	if !isUserInConversation {
		return Message{}, fmt.Errorf(ErrUserNotInConversation)
	}

	// Check if the user is in the forwoarded conversation
	isUserInConversation, err = db.IsUserInConversation(id_conversation, id_auth)
	if err != nil {
		return Message{}, fmt.Errorf("database error checking if user is in conversation: %w", err)
	}
	if !isUserInConversation {
		return Message{}, fmt.Errorf(ErrUserNotInConversation)
	}

	message.Forward = message.Author
	message.Author = id_auth
	message.Recipient = id_conversation

	res, err := db.c.Exec(`
	INSERT INTO messages (
	text,
	photo, 
	author, 
	recipient,
	forward
	) VALUES ($1, $2, $3, $4, $5);`, message.Text, message.Photo, message.Author, message.Recipient, message.Forward)
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

func (db *appdbimpl) ReactMessage(id_message int64, id_auth int64, reaction Reaction) error {

	if reaction.Reaction == "" {
		return fmt.Errorf("empty reaction")
	}

	// Check if the message existss
	message, err := db.GetMessage(id_message)
	if err != nil && err.Error() == ErrMessageNotFound {
		return err
	}
	if err != nil {
		return fmt.Errorf("error getting message: %w", err)
	}

	var id_conversation int64 = message.Recipient

	// Check if the user is in the conversation
	isUserInConversation, err := db.IsUserInConversation(id_conversation, id_auth)
	if err != nil {
		return fmt.Errorf("database error checking if user is in conversation: %w", err)
	}
	if !isUserInConversation {
		return fmt.Errorf(ErrUserNotInConversation)
	}

	// Check if the user has already reacted to the message
	hasReaction, err := db.HasReaction(id_message, id_auth)
	if err != nil {
		return fmt.Errorf("database error checking if user has already reacted to message: %w", err)
	}

	if hasReaction && reaction.Reaction != "" {
		_, err = db.c.Exec(`
		UPDATE reactions
		SET
			reaction = $1
		WHERE
			id_user = $2
		AND
			id_message = $3;

		`, reaction.Reaction, id_auth, id_message)
		if err != nil {
			return fmt.Errorf("database error reacting message: %w", err)
		}

		return nil
	}

	// React Mesasge
	_, err = db.c.Exec(`
	INSERT INTO reactions (
		id_user,
		id_message,
		reaction
	) VALUES ($1, $2, $3);`, id_auth, id_message, reaction.Reaction)
	if err != nil {
		return fmt.Errorf("database error reacting message: %w", err)
	}

	return nil
}

func (db *appdbimpl) UnReactMessage(id_message int64, id_auth int64) error {
	// Check if the message existss
	message, err := db.GetMessage(id_message)
	if err != nil && err.Error() == ErrMessageNotFound {
		return err
	}
	if err != nil {
		return fmt.Errorf("error getting message: %w", err)
	}

	var id_conversation int64 = message.Recipient

	// Check if the user is in the conversation
	isUserInConversation, err := db.IsUserInConversation(id_conversation, id_auth)
	if err != nil {
		return fmt.Errorf("database error checking if user is in conversation: %w", err)
	}
	if !isUserInConversation {
		return fmt.Errorf(ErrUserNotInConversation)
	}

	// Check if the user has already reacted to the message
	hasReaction, err := db.HasReaction(id_message, id_auth)
	if err != nil {
		return fmt.Errorf("database error checking if user has already reacted to message: %w", err)
	}

	// Remove reaction
	if hasReaction {
		_, err = db.c.Exec(`
		DELETE FROM reactions
		WHERE
			id_user = $1
		AND
			id_message = $2;

		`, id_auth, id_message)
		if err != nil {
			return fmt.Errorf("database error reacting message: %w", err)
		}

		return nil
	}

	return nil
}
