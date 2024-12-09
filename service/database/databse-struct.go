package database

type User struct {
	Id    int64  `json:"id"`
	Name  string `json:"name"`
	Photo []byte `json:"photo"`
}

// type Message struct {
// 	Id      int
// 	content string
// 	photo   []byte

// 	author       User
// 	recipient    Conversation
// 	forwarded_to User

// 	timestamp string
// 	reaction  Reaction
// }

// type Reaction struct {
// 	user     User
// 	reaction string
// }

// type Conversation struct {
// 	Id           int
// 	name         string
// 	photo        []byte
// 	cnv_type     string
// 	participants []User
// 	messages     []Message
// }
