package database

import (
	"fmt"
)

// func (db *appdbimpl) loginUser(name string) (User, error) {
// 	var user User

// 	err := db.c.QueryRow(`SELECT id, name FROM users WHERE name = $1`, name).Scan(&user.Id, &user.Name)

// 	// if user not present in database, add new user
// 	if err == sql.ErrNoRows {
// 		id, err := db.AddUser(User{Name: name})
// 		if err != nil {
// 			return user, fmt.Errorf("error adding user: %w", err)
// 		}
// 		user.Id = id
// 		return user, nil
// 	}
// 	if err != nil {
// 		return user, fmt.Errorf("error getting all users: %w", err)
// 	}

// 	return user, nil
// }

func (db *appdbimpl) GetUsers() ([]User, error) {
	var users []User
	rows, err := db.c.Query(`SELECT id, name FROM users`)
	if err != nil {
		return users, fmt.Errorf("error getting all users: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var user User
		if err := rows.Scan(&user.Id, &user.Name); err != nil {
			return users, fmt.Errorf("error getting all users: %w", err)
		}
		users = append(users, user)
	}
	if err := rows.Err(); err != nil {
		return users, fmt.Errorf("error unpacking all users: %w", err)
	}

	return users, nil
}

// func (db *appdbimpl) GetUser(id int) (User, error) {
// 	var user User

// 	row := db.c.QueryRow(`SELECT id, name FROM users WHERE id = $1`, id)
// 	if err := row.Scan(&user.Id, &user.Name); err != nil {
// 		return user, fmt.Errorf("error getting user: %w", err)
// 	}

// 	return user, nil
// }

// func (db *appdbimpl) AddUser(user User) (int64, error) {
// 	var id int64
// 	var err error

// 	res, err := db.c.Exec(`INSERT INTO users (name) VALUES ($1)`, user.Name)
// 	if err != nil {
// 		return id, fmt.Errorf("error adding user: %w", err)
// 	}

// 	id, err = res.LastInsertId()

// 	return id, nil
// }
