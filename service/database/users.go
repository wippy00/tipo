package database

import (
	"database/sql"
	"errors"
	"fmt"
)

func (db *appdbimpl) UserExistById(id int64) (User, bool, error) {
	var user User

	err := db.c.QueryRow(`SELECT id, name, photo FROM users WHERE id = $1`, id).Scan(&user.Id, &user.Name, &user.Photo)

	if errors.Is(err, sql.ErrNoRows) {
		return user, false, nil
	}
	if err != nil {
		return user, false, fmt.Errorf("error getting user: %w", err)
	}

	return user, true, nil

}

func (db *appdbimpl) UserExistByName(name string) (User, bool, error) {
	var user User

	err := db.c.QueryRow(`SELECT id, name, photo FROM users WHERE name = $1`, name).Scan(&user.Id, &user.Name, &user.Photo)

	if errors.Is(err, sql.ErrNoRows) {
		return user, false, nil
	}
	if err != nil {
		return user, false, fmt.Errorf("error getting user: %w", err)
	}

	return user, true, nil

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
	err = rows.Err()
	if err != nil {
		return users, fmt.Errorf("error unpacking all users: %w", err)
	}

	return users, nil
}

func (db *appdbimpl) GetUser(id int64) (User, error) {
	var user User

	row := db.c.QueryRow(`SELECT id, name, photo FROM users WHERE id = $1`, id)
	if err := row.Scan(&user.Id, &user.Name, &user.Photo); err != nil {
		return user, fmt.Errorf("User not exist: %w", err)
	}

	return user, nil
}

func (db *appdbimpl) AddUser(user User) (User, error) {
	var err error

	res, err := db.c.Exec(`INSERT INTO users (name) VALUES ($1)`, user.Name)
	if err != nil {
		return user, fmt.Errorf("error adding user: %w", err)
	}

	user.Id, err = res.LastInsertId()
	if err != nil {
		return user, fmt.Errorf("error getting last inserted id of added user: %w", err)
	}

	return user, nil
}

func (db *appdbimpl) LoginUser(name string) (user User, isNew bool, err error) {

	// check if user exists
	user, exist, err := db.UserExistByName(name)
	if err != nil {
		return user, false, fmt.Errorf("error checking user existence: %w", err)
	}

	// if exist return user	and isNew = false
	if exist {
		return user, false, nil
	}

	// if not exist return added user	and isNew = true
	added_user, err := db.AddUser(User{Name: name})
	if err != nil {
		return added_user, false, fmt.Errorf("error adding user: %w", err)
	}

	return added_user, true, nil
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
