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

func DbUser(user User) database.User {
	return database.User{
		Id:    user.Id,
		Name:  user.Name,
		Photo: user.Photo,
	}
}

type Conversation struct {
	Id           int64   `json:"id"`
	Name         string  `json:"name"`
	Photo        []byte  `json:"photo"`
	Cnv_type     string  `json:"cnv_type"`
	Participants []User  `json:"participants"`
	Last_message Message `json:"last_message"`
}

func convertParticipants(dbParticipants []database.User) []User {
	participants := make([]User, len(dbParticipants))
	for i, dbUser := range dbParticipants {
		participants[i] = NewUser(dbUser)
	}
	return participants
}
func convertDbParticipants(participants []User) []database.User {
	dbParticipants := make([]database.User, len(participants))
	for i, user := range participants {
		dbParticipants[i] = DbUser(user)
	}
	return dbParticipants

}

func NewConversation(conversation database.Conversation) Conversation {
	return Conversation{
		Id:           conversation.Id,
		Name:         conversation.Name,
		Photo:        conversation.Photo,
		Cnv_type:     conversation.Cnv_type,
		Participants: convertParticipants(conversation.Participants),
		Last_message: NewMessage(conversation.Last_message),
	}
}

func DbConversation(conversation Conversation) database.Conversation {
	return database.Conversation{
		Id:           conversation.Id,
		Name:         conversation.Name,
		Photo:        conversation.Photo,
		Cnv_type:     conversation.Cnv_type,
		Participants: convertDbParticipants(conversation.Participants),
		Last_message: DbMessage(conversation.Last_message),
	}
}

type Message struct {
	Id               int64     `json:"id"`
	Text             string    `json:"text"`
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
		Text:             message.Text,
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
		Text:             message.Text,
		Photo:            message.Photo,
		Author:           message.Author,
		Recipient:        message.Recipient,
		Forwarded_source: message.Forwarded_source,
		Timestamp:        message.Timestamp,
	}
}
