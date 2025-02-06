/*
Package database is the middleware between the app database and the code. All data (de)serialization (save/load) from a
persistent database are handled here. Database specific logic should never escape this package.

To use this package you need to apply migrations to the database if needed/wanted, connect to it (using the database
data source name from config), and then initialize an instance of AppDatabase from the DB connection.

For example, this code adds a parameter in `webapi` executable for the database data source name (add it to the
main.WebAPIConfiguration structure):

	DB struct {
		Filename string `conf:""`
	}

This is an example on how to migrate the DB and connect to it:

	// Start Database
	logger.Println("initializing database support")
	db, err := sql.Open("sqlite3", "./foo.db")
	if err != nil {
		logger.WithError(err).Error("error opening SQLite DB")
		return fmt.Errorf("opening SQLite: %w", err)
	}
	defer func() {
		logger.Debug("database stopping")
		_ = db.Close()
	}()

Then you can initialize the AppDatabase and pass it to the api package.
*/
package database

import (
	"database/sql"
	"errors"
	"fmt"
)

// AppDatabase is the high level interface for the DB
type AppDatabase interface {
	GetName() (string, error)
	SetName(name string) error

	// User
	LoginUser(name string) (User, bool, error)
	UpdateUserName(id int64, name string) (User, error)
	UpdateUserPhoto(id int64, photo []byte) (User, error)
	UserExistById(id int64) (User, bool, error)
	GetUsers() ([]User, error)
	GetUser(id int64) (User, error)

	// Conversation
	GetConversation(id int64) (Conversation, error)
	GetConversationsOfUser(id int64) ([]Conversation, error)

	// GetConversationMessages(id int64)
	UpdateConversationName(id int64, id_auth int64, new_name string) (Conversation, error)
	UpdateConversationPhoto(id int64, id_auth int64, new_photo []byte) (Conversation, error)
	AddUserToConversation(id_conversation int64, id_auth int64, id_user int64) (Conversation, error)
	RemoveUserFromConversation(id_conversation int64, id_auth int64, id_user int64) error

	CreateConversation(id_auth int64, conversation Conversation) (Conversation, error)

	// Message
	GetMessagesOfConversation(id_conversation int64, id_auth int64) ([]Message, error)

	SendMessage(id_conversation int64, id_auth int64, message Message) (Message, error)
	DeleteMessage(id_message int64, id_auth int64) error
	ForwardMessage(id_message int64, id_auth int64, id_conversation int64) (Message, error)
	ReactMessage(id_message int64, id_auth int64, reaction Reaction) error
	UnReactMessage(id_message int64, id_auth int64) error

	Ping() error
}

type appdbimpl struct {
	c *sql.DB
}

// New returns a new instance of AppDatabase based on the SQLite connection `db`.
// `db` is required - an error will be returned if `db` is `nil`.
func New(db *sql.DB) (AppDatabase, error) {
	if db == nil {
		return nil, errors.New("database is required when building a AppDatabase")
	}

	// Check if table existSs. If not, the database is empty, and we need to create the structure
	var tableName string
	err := db.QueryRow(`SELECT name FROM sqlite_master WHERE type='table' AND name='users';`).Scan(&tableName)
	if errors.Is(err, sql.ErrNoRows) {
		// ############################################################
		// ###						users							###
		// ############################################################
		user_table := `CREATE TABLE IF NOT EXISTS users (
			id INTEGER NOT NULL PRIMARY KEY, 
			name VARCHAR(30) NOT NULL UNIQUE,
			photo BLOB
		);`
		_, err = db.Exec(user_table)
		if err != nil {
			return nil, fmt.Errorf("error creating database structure user: %w", err)
		}
		// ############################################################
		// ###					conversations						###
		// ############################################################
		conversations_table := `CREATE TABLE IF NOT EXISTS conversations (
			id INTEGER NOT NULL PRIMARY KEY, 
			name VARCHAR(50) NOT NULL,
			photo BLOB,
			cnv_type VARCHAR(10) NOT NULL
		);`
		_, err = db.Exec(conversations_table)
		if err != nil {
			return nil, fmt.Errorf("error creating database structure conversations: %w", err)
		}
		// ############################################################
		// ###				conversations_members					###
		// ############################################################
		conversations_members_table := `CREATE TABLE IF NOT EXISTS conversations_members (
			id_conversations INTEGER NOT NULL, 
			id_user INTEGER NOT NULL,
	
			FOREIGN KEY (id_conversations) REFERENCES conversations(id) ON DELETE CASCADE,
			FOREIGN KEY (id_user) REFERENCES users(id) ON DELETE CASCADE
		);`
		_, err = db.Exec(conversations_members_table)
		if err != nil {
			return nil, fmt.Errorf("error creating database structure conversations_members: %w", err)
		}
		// ############################################################
		// ###						messages						###
		// ############################################################
		messages_table := `CREATE TABLE IF NOT EXISTS messages (
			id INTEGER NOT NULL PRIMARY KEY, 
			text TEXT,
			photo BLOB,
			author INT NOT NULL,
			recipient INT NOT NULL,
			forward INT,
			reply INT,
			timestamp TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			
			FOREIGN KEY (author) REFERENCES users(id) ON DELETE CASCADE,
			FOREIGN KEY (recipient) REFERENCES conversations(id) ON DELETE CASCADE,
			FOREIGN KEY (forward) REFERENCES users(id) ON DELETE CASCADE,
			FOREIGN KEY (reply) REFERENCES messages(id) ON DELETE CASCADE
		);`
		_, err = db.Exec(messages_table)
		if err != nil {
			return nil, fmt.Errorf("error creating database structure messages: %w", err)
		}
		// ############################################################
		// ###					messages_readers					###
		// ############################################################
		messages_readers_table := `CREATE TABLE IF NOT EXISTS messages_readers (
			id_user INTEGER NOT NULL,
			id_message INTEGER NOT NULL,
			
			FOREIGN KEY (id_user) REFERENCES users(id) ON DELETE CASCADE,
			FOREIGN KEY (id_message) REFERENCES messages(id) ON DELETE CASCADE
		);`
		_, err = db.Exec(messages_readers_table)
		if err != nil {
			return nil, fmt.Errorf("error creating database structure messages_readers: %w", err)
		}
		// ############################################################
		// ###						reactions						###
		// ############################################################
		reactions_table := `CREATE TABLE IF NOT EXISTS reactions (
			id_user INTEGER NOT NULL,
			id_message INTEGER NOT NULL,
			reaction VARCHAR(3) NOT NULL,
			
			FOREIGN KEY (id_user) REFERENCES users(id) ON DELETE CASCADE,
			FOREIGN KEY (id_message) REFERENCES messages(id) ON DELETE CASCADE
		);`
		_, err = db.Exec(reactions_table)
		if err != nil {
			return nil, fmt.Errorf("error creating database structure reactions: %w", err)
		}
	}

	return &appdbimpl{
		c: db,
	}, nil
}

func (db *appdbimpl) Ping() error {
	return db.c.Ping()
}
