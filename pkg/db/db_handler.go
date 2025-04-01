package db

import (
	"database/sql"
	"errors"
	"strings"

	_ "github.com/lib/pq"
)

func userExist(db *sql.DB, username string) (int, error) {
	username = strings.ToLower(username)
	var userID int
	err := db.QueryRow(`SELECT id FROM users WHERE username = $1`, username).Scan(&userID)
	if err == sql.ErrNoRows {
		return -1, nil
	} else if err != nil {
		return -1, err
	}
	return userID, nil
}

func AddUser(db *sql.DB, d UserData) (int, error) {
	d.UserName = strings.ToLower(d.UserName)
	userID, err := userExist(db, d.UserName)
	if err != nil {
		return -1, err
	}
	if userID != -1 {
		return -1, errors.New("user already exists")
	}

	insertStatement := `INSERT INTO users (username) VALUES ($1) RETURNING id`
	err = db.QueryRow(insertStatement, d.UserName).Scan(&userID)
	if err != nil {
		return -1, err
	}

	insertStatement = `INSERT INTO userdata (userid, name, surname, description) VALUES ($1, $2, $3, $4)`
	_, err = db.Exec(insertStatement, userID, d.Name, d.Surname, d.Description)
	if err != nil {
		return -1, err
	}

	return userID, nil
}

func DeleteUser(db *sql.DB, id int) error {
	_, err := db.Exec(`DELETE FROM userdata WHERE userid=$1`, id)
	if err != nil {
		return err
	}

	_, err = db.Exec(`DELETE FROM users WHERE id=$1`, id)
	return err
}

func ListUsers(db *sql.DB) ([]UserData, error) {
	var users []UserData
	query := `SELECT u.id, u.username, ud.name, ud.surname, ud.description 
			FROM users u 
			JOIN userdata ud ON u.id = ud.userid`
	rows, err := db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var u UserData
		if err := rows.Scan(&u.ID, &u.UserName, &u.Name, &u.Surname, &u.Description); err != nil {
			return nil, err
		}
		users = append(users, u)
	}
	return users, nil
}

func UpdateUser(db *sql.DB, d UserData) error {
	d.UserName = strings.ToLower(d.UserName)
	updateStatement := `UPDATE userdata SET name=$1, surname=$2, description=$3 WHERE userid=$4`
	_, err := db.Exec(updateStatement, d.Name, d.Surname, d.Description, d.ID)
	return err
}

func GetUser(db *sql.DB, id int) (*UserData, error) {
	var user UserData
	query := `SELECT u.id, u.username, ud.name, ud.surname, ud.description 
			FROM users u 
			JOIN userdata ud ON u.id = ud.userid 
			WHERE u.id = $1`
	err := db.QueryRow(query, id).Scan(&user.ID, &user.UserName, &user.Name, &user.Surname, &user.Description)
	if err == sql.ErrNoRows {
		return nil, errors.New("user not found")
	} else if err != nil {
		return nil, err
	}
	return &user, nil
}
