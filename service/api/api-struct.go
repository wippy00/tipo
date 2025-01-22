package api

import (
	"time"

	"github.com/wippy00/wasa-text/service/database"
)

type User struct {
	Id    int64  `json:"id"`
	Name  string `json:"name"`
	Photo []byte `json:"photo"`
}

func NewUser(user database.User) User {
	return User{
		Id:    user.Id,
		Name:  user.Name,
		Photo: user.Photo,
	}
}

type Conversation struct {
	Id           int64  `json:"id"`
	Name         string `json:"name"`
	Photo        []byte `json:"photo"`
	CnvType      string `json:"cnvType"`
	Participants []User `json:"participants"`
	// messages     []Message `json:"messages"`
}

func convertParticipants(dbParticipants []database.User) []User {
	participants := make([]User, len(dbParticipants))
	for i, dbUser := range dbParticipants {
		participants[i] = NewUser(dbUser)
	}
	return participants
}

func NewConversation(conversation database.Conversation) Conversation {
	return Conversation{
		Id:           conversation.Id,
		Name:         conversation.Name,
		Photo:        conversation.Photo,
		CnvType:      conversation.CnvType,
		Participants: convertParticipants(conversation.Participants),
	}
}

type Message struct {
	Id               int64     `json:"id"`
	Content          string    `json:"content"`
	Photo            []byte    `json:"photo"`
	Author           int64     `json:"author"`
	Recipient        int64     `json:"recipient"`
	Forwarded_source int64     `json:"forwarded_source"`
	Timestamp        time.Time `json:"timestamp"`
	// reaction  Reaction
}

func NewMessage(message database.Message) Message {
	return Message{
		Id:               message.Id,
		Content:          message.Content,
		Photo:            message.Photo,
		Author:           message.Author,
		Recipient:        message.Recipient,
		Forwarded_source: message.Forwarded_source,
		Timestamp:        message.Timestamp,
	}
}

func DbMessage(message Message) database.Message {
	return database.Message{
		Id:               message.Id,
		Content:          message.Content,
		Photo:            message.Photo,
		Author:           message.Author,
		Recipient:        message.Recipient,
		Forwarded_source: message.Forwarded_source,
		Timestamp:        message.Timestamp,
	}
}
