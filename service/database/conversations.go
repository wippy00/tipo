package database

import (
	"fmt"
)

func (db *appdbimpl) IsGroup(id_conversation int64) (bool, error) {
	var cnv_type string

	err := db.c.QueryRow(`SELECT cnv_type FROM conversations WHERE id = $1`, id_conversation).Scan(&cnv_type)
	if err != nil {
		return false, fmt.Errorf("error getting conversation type: %w", err)
	}

	if cnv_type == "group" {
		return true, nil
	} else {
		return false, nil
	}
}

func (db *appdbimpl) IsUserInGroup(id_conversation int64, id_user int64) (bool, error) {
	var rowsCount int64

	err := db.c.QueryRow(`SELECT count(*) AS rowsCount FROM conversations_members WHERE id_conversations = $1 AND id_user = $2`, id_conversation, id_user).Scan(&rowsCount)
	if err != nil {
		return false, fmt.Errorf("database error checking if user is in group: %w", err)
	}
	if rowsCount > 0 {
		return true, nil
	} else {
		return false, nil
	}

}

// ############################################################

func (db *appdbimpl) GetConversation(id int64) (Conversation, error) {
	var conversations []Conversation

	rows, err := db.c.Query(`
	SELECT 
		conversations.id as conversation_id,
		conversations.name as conversation_name,
		conversations.photo as conversation_photo,
		conversations.cnv_type as conversation_type,
		participants.id AS participant_id,
		participants.name AS participant_name,
		participants.photo AS participant_photo
	FROM 
		users
	JOIN 
		conversations_members ON users.id = conversations_members.id_user
	JOIN 
		conversations ON conversations_members.id_conversations = conversations.id
	JOIN 
		conversations_members conversations_members_2 ON conversations.id = conversations_members_2.id_conversations
	JOIN 
		users participants ON conversations_members_2.id_user = participants.id
	WHERE 
		conversations.id = $1
	GROUP BY
        participants.id
	ORDER BY
		participants.id; 
	`, id)
	if err != nil {
		return Conversation{}, fmt.Errorf("error getting conversation: %w", err)
	}
	defer rows.Close()

	conversation_map := make(map[int64]*Conversation)
	lastParticipantID := int64(-1)
	for rows.Next() {
		var conversation Conversation
		var participant User

		if err := rows.Scan(
			&conversation.Id,
			&conversation.Name,
			&conversation.Photo,
			&conversation.CnvType,
			&participant.Id,
			&participant.Name,
			&participant.Photo,
		); err != nil {
			return Conversation{}, fmt.Errorf("error getting conversation row: %w", err)
		}

		// Verifica e crea la conversazione se non esiste
		if conversation_map[conversation.Id] == nil {
			conversation_map[conversation.Id] = &Conversation{
				Id:      conversation.Id,
				Name:    conversation.Name,
				Photo:   conversation.Photo,
				CnvType: conversation.CnvType,
			}
			lastParticipantID = -1
		}

		// Aggiungi il partecipante
		if lastParticipantID != participant.Id {
			conversation_map[conversation.Id].Participants = append(
				conversation_map[conversation.Id].Participants,
				User{
					Id:    participant.Id,
					Name:  participant.Name,
					Photo: participant.Photo,
				},
			)
		}
		lastParticipantID = participant.Id

	}

	conversations = make([]Conversation, 0, len(conversation_map))
	for _, conv := range conversation_map {
		conversations = append(conversations, *conv)
	}

	if len(conversations) > 0 {
		return conversations[0], nil
	} else {
		return Conversation{}, fmt.Errorf("no conversation found with id %d", id)
	}
}

func (db *appdbimpl) GetConversationsOfUser(id int64) ([]Conversation, error) {
	var conversations []Conversation

	rows, err := db.c.Query(`
	SELECT 
		conversations.id as conversation_id,
		conversations.name as conversation_name,
		conversations.photo as conversation_photo,
		conversations.cnv_type as conversation_type,
		participants.id AS participant_id,
		participants.name AS participant_name,
		participants.photo AS participant_photo
	FROM 
		users
	JOIN 
		conversations_members ON users.id = conversations_members.id_user
	JOIN 
		conversations ON conversations_members.id_conversations = conversations.id
	JOIN 
		conversations_members conversations_members_2 ON conversations.id = conversations_members_2.id_conversations
	JOIN 
		users participants ON conversations_members_2.id_user = participants.id
	WHERE 
		users.id = $1
	ORDER BY
		participants.id; 
	`, id)
	if err != nil {
		return conversations, fmt.Errorf("error getting conversation: %w", err)
	}
	defer rows.Close()

	conversation_map := make(map[int64]*Conversation)
	lastParticipantID := int64(-1)
	for rows.Next() {
		var conversation Conversation
		var participant User

		if err := rows.Scan(
			&conversation.Id,
			&conversation.Name,
			&conversation.Photo,
			&conversation.CnvType,
			&participant.Id,
			&participant.Name,
			&participant.Photo,
		); err != nil {
			return conversations, fmt.Errorf("error getting conversation row: %w", err)
		}

		// Verifica e crea la conversazione se non esiste
		if conversation_map[conversation.Id] == nil {
			conversation_map[conversation.Id] = &Conversation{
				Id:      conversation.Id,
				Name:    conversation.Name,
				Photo:   conversation.Photo,
				CnvType: conversation.CnvType,
			}
			lastParticipantID = -1
		}

		// Aggiungi il partecipante
		if lastParticipantID != participant.Id {
			conversation_map[conversation.Id].Participants = append(
				conversation_map[conversation.Id].Participants,
				User{
					Id:    participant.Id,
					Name:  participant.Name,
					Photo: participant.Photo,
				},
			)
		}
		lastParticipantID = participant.Id

	}

	conversations = make([]Conversation, 0, len(conversation_map))
	for _, conv := range conversation_map {
		conversations = append(conversations, *conv)
	}

	return conversations, nil
}

func (db *appdbimpl) UpdateConversationName(id_conversation int64, id_auth int64, new_name string) (Conversation, error) {
	var conversation Conversation

	//Check if is a group
	isGroup, err := db.IsGroup(id_conversation)
	if err != nil {
		return Conversation{}, fmt.Errorf("database error checking if is group: %w", err)
	}
	if !isGroup {
		return Conversation{}, fmt.Errorf("conversation is not a group")
	}

	//Check if user is in group
	isUserInGroup, err := db.IsUserInGroup(id_conversation, id_auth)
	if err != nil {
		return Conversation{}, fmt.Errorf("database error checking if user is in group: %w", err)
	}
	if !isUserInGroup {
		return Conversation{}, fmt.Errorf("auth user is not in this group")
	}

	//Update conversation name
	res, err := db.c.Exec(`UPDATE conversations SET name = $1 WHERE id = $2`, new_name, id_conversation)
	if err != nil {
		return Conversation{}, fmt.Errorf("database error updating conversation name: %w", err)
	}
	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return Conversation{}, fmt.Errorf("database error getting rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return Conversation{}, fmt.Errorf("no conversation with id %d found", id_conversation)
	}

	conversation, err = db.GetConversation(id_conversation)
	if err != nil {
		return Conversation{}, fmt.Errorf("database error getting conversation: %w", err)
	}

	return conversation, nil
}

func (db *appdbimpl) UpdateConversationPhoto(id_conversation int64, id_auth int64, new_photo []byte) (Conversation, error) {
	var conversation Conversation

	//Check if is a group
	isGroup, err := db.IsGroup(id_conversation)
	if err != nil {
		return Conversation{}, fmt.Errorf("database error checking if is group: %w", err)
	}
	if !isGroup {
		return Conversation{}, fmt.Errorf("conversation is not a group")
	}

	//Check if user is in group
	isUserInGroup, err := db.IsUserInGroup(id_conversation, id_auth)
	if err != nil {
		return Conversation{}, fmt.Errorf("database error checking if user is in group: %w", err)
	}
	if !isUserInGroup {
		return Conversation{}, fmt.Errorf("auth user is not in this group")
	}

	res, err := db.c.Exec(`UPDATE conversations SET photo = $1 WHERE id = $2`, new_photo, id_conversation)
	if err != nil {
		return Conversation{}, fmt.Errorf("database error updating conversation photo: %w", err)
	}
	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return Conversation{}, fmt.Errorf("database error getting rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return Conversation{}, fmt.Errorf("no conversation with id %d found", id_conversation)
	}

	conversation, err = db.GetConversation(id_conversation)
	if err != nil {
		return Conversation{}, fmt.Errorf("database error getting conversation: %w", err)
	}

	return conversation, nil
}

func (db *appdbimpl) AddUserToConversation(id_conversation int64, id_auth int64, id_user int64) (Conversation, error) {

	//Check if is a group
	isGroup, err := db.IsGroup(id_conversation)
	if err != nil {
		return Conversation{}, fmt.Errorf("database error checking if is group: %w", err)
	}
	if !isGroup {
		return Conversation{}, fmt.Errorf("conversation is not a group")
	}

	//Check if user is in group
	isUserInGroup, err := db.IsUserInGroup(id_conversation, id_auth)
	if err != nil {
		return Conversation{}, fmt.Errorf("database error checking if user is in group: %w", err)
	}
	if !isUserInGroup {
		return Conversation{}, fmt.Errorf("auth user is not in this group")
	}

	//Check if user is already in conversation
	var rowsCount int64
	erri := db.c.QueryRow(`SELECT count(*) AS rowsCount FROM conversations_members WHERE id_conversations = $1 AND id_user = $2`, id_conversation, id_user).Scan(&rowsCount)
	if erri != nil {
		return Conversation{}, fmt.Errorf("database error checking if user is already in conversation: %w", err)
	}
	if rowsCount > 0 {
		return Conversation{}, fmt.Errorf("user is already in conversation")
	}

	// Add user to conversation
	_, erro := db.c.Exec(`INSERT INTO conversations_members (id_conversations, id_user) VALUES ($1, $2)`, id_conversation, id_user)
	if erro != nil {
		return Conversation{}, fmt.Errorf("database error adding user to conversation: %w", err)
	}

	conversation, err := db.GetConversation(id_conversation)
	if err != nil {
		return conversation, fmt.Errorf("database error getting conversation: %w", err)
	}

	return conversation, nil
}
