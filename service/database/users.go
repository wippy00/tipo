package database

import (
	"database/sql"
	"fmt"
)

func (db *appdbimpl) LoginUser(name string) (User, error) {
	var user User

	err := db.c.QueryRow(`SELECT id, name FROM users WHERE name = $1`, name).Scan(&user.Id, &user.Name)

	// if user not present in database, add new user
	if err == sql.ErrNoRows {
		id, err := db.AddUser(User{Name: name})
		if err != nil {
			return user, fmt.Errorf("error adding user: %w", err)
		}
		user.Id = id
		user.Name = name
		return user, nil
	}
	if err != nil {
		return user, fmt.Errorf("error getting all users: %w", err)
	}

	return user, nil
}

func (db *appdbimpl) GetUsers() ([]User, error) {
	var users []User
	rows, err := db.c.Query(`SELECT id, name, photo FROM users`)
	if err != nil {
		return users, fmt.Errorf("error getting all users: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var user User
		if err := rows.Scan(&user.Id, &user.Name, &user.Photo); err != nil {
			return users, fmt.Errorf("error getting all users: %w", err)
		}
		users = append(users, user)
	}
	if err := rows.Err(); err != nil {
		return users, fmt.Errorf("error unpacking all users: %w", err)
	}

	return users, nil
}

func (db *appdbimpl) GetUser(id int64) (User, error) {
	var user User

	row := db.c.QueryRow(`SELECT id, name, photo FROM users WHERE id = $1`, id)
	if err := row.Scan(&user.Id, &user.Name, &user.Photo); err != nil {
		return user, fmt.Errorf("error getting user: %w", err)
	}

	return user, nil
}

func (db *appdbimpl) AddUser(user User) (int64, error) {
	var id int64
	var err error

	res, err := db.c.Exec(`INSERT INTO users (name) VALUES ($1)`, user.Name)
	if err != nil {
		return id, fmt.Errorf("error adding user: %w", err)
	}

	id, err = res.LastInsertId()
	if err != nil {
		return id, fmt.Errorf("error getting last inserted id of added user: %w", err)
	}

	return id, nil
}

func (db *appdbimpl) UpdateUserName(id int64, new_name string) (User, error) {
	var user User

	res, err := db.c.Exec(`UPDATE users SET name = $1 WHERE id = $2`, new_name, id)
	if err != nil {
		return user, fmt.Errorf("database error updating user name: %w", err)
	}
	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return user, fmt.Errorf("database error getting rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return user, fmt.Errorf("no user with id %d found", id)
	}

	user, err = db.GetUser(id)
	if err != nil {
		return user, fmt.Errorf("database error getting user: %w", err)
	}

	return user, nil
}

func (db *appdbimpl) UpdateUserPhoto(id int64, new_photo []byte) (User, error) {
	var user User

	res, err := db.c.Exec(`UPDATE users SET photo = $1 WHERE id = $2`, new_photo, id)
	if err != nil {
		return user, fmt.Errorf("database error updating user photo: %w", err)
	}
	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return user, fmt.Errorf("database error getting rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return user, fmt.Errorf("no user with id %d found", id)
	}

	user, err = db.GetUser(id)
	if err != nil {
		return user, fmt.Errorf("database error getting user: %w", err)
	}

	return user, nil
}
