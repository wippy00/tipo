package database

import "time"

type User struct {
	Id    int64  `json:"id"`
	Name  string `json:"name"`
	Photo []byte `json:"photo"`
}

type Conversation struct {
	Id           int64   `json:"id"`
	Name         string  `json:"name"`
	Photo        []byte  `json:"photo"`
	Cnv_type     string  `json:"cnv_type"`
	Participants []User  `json:"participants"`
	Last_message Message `json:"last_message"`
}

type Message struct {
	Id        int64      `json:"id"`
	Text      string     `json:"text"`
	Photo     []byte     `json:"photo"`
	Author    int64      `json:"author"`
	Recipient int64      `json:"recipient"`
	Forward   int64      `json:"forward"`
	Reply     int64      `json:"reply"`
	Timestamp time.Time  `json:"timestamp"`
	Reactions []Reaction `json:"reactions"`
	Read      bool       `json:"read"`
}

type Reaction struct {
	User     int64  `json:"user"`
	Reaction string `json:"reaction"`
}
