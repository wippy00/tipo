package database

import (
	"fmt"
)

func (db *appdbimpl) GetConversation(id int64) ([]Conversation, error) {
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

func (db *appdbimpl) GetConversationOfUser(id int64) ([]Conversation, error) {
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
