package database

import (
	"fmt"
)

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
