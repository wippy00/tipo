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
	LoginUser(name string) (User, error)

	GetUsers() ([]User, error)
	// GetUser(id int64) (User, error)
	// AddUser(user User) error

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
		// ###						user							###
		// ############################################################
		user_table := `CREATE TABLE IF NOT EXISTS users (
			id INTEGER NOT NULL PRIMARY KEY, 
			name VARCHAR(30) NOT NULL,
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
			cnv_type TEXT NOT NULL
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
	
			FOREIGN KEY (id_conversations) REFERENCES conversations(id),
			FOREIGN KEY (id_user) REFERENCES users(id)
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
			forwarded_to INT,
			timestamp TEXT DEFAULT CURRENT_TIMESTAMP,
			
			FOREIGN KEY (author) REFERENCES users(id),
			FOREIGN KEY (recipient) REFERENCES conversations(id),
			FOREIGN KEY (forwarded_to) REFERENCES conversations(id)
		);`
		_, err = db.Exec(messages_table)
		if err != nil {
			return nil, fmt.Errorf("error creating database structure messages: %w", err)
		}
		// ############################################################
		// ###					messages_readers					###
		// ############################################################
		messages_readers_table := `CREATE TABLE IF NOT EXISTS messages_readers (
			id_user INTEGER NOT NULL PRIMARY KEY,
			id_message INTEGER NOT NULL,
			
			FOREIGN KEY (id_user) REFERENCES users(id),
			FOREIGN KEY (id_message) REFERENCES messages(id)
		);`
		_, err = db.Exec(messages_readers_table)
		if err != nil {
			return nil, fmt.Errorf("error creating database structure messages_readers: %w", err)
		}
		// ############################################################
		// ###						reactions						###
		// ############################################################
		reactions_table := `CREATE TABLE IF NOT EXISTS reactions (
			id_user INTEGER NOT NULL PRIMARY KEY,
			id_message INTEGER NOT NULL,
			reaction VARCHAR(3) NOT NULL,
			
			FOREIGN KEY (id_user) REFERENCES users(id),
			FOREIGN KEY (id_message) REFERENCES messages(id)
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
