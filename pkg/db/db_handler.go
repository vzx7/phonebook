package db

import (
	"fmt"
	"strings"

	_ "github.com/lib/pq"
)

func userExist(_username string) int {
	username := strings.ToLower(_username)

	db, err := OpenConnection([]string{})
	if err != nil {
		fmt.Println(err)
		return -1
	}
	defer db.Close()

	userID := -1
	statement := fmt.Sprintf(`SELECT "id" FROM "users" WHERE username = %s`, username)
	rows, err := db.Query(statement)
	if err != nil {
		fmt.Println(err)
		return -1
	}

	for rows.Next() {
		var id int
		err = rows.Scan(&id)
		if err != nil {
			fmt.Println("Scan", err)
			return -1
		}
		userID = id
	}

	defer rows.Close()
	return userID
}

func AddUser(d UserData) int {
	d.UserName = strings.ToLower(d.UserName)

	db, err := OpenConnection([]string{})
	if err != nil {
		fmt.Println(err)
		return -1
	}
	defer db.Close()

	userID := userExist(d.UserName)
	if userID != -1 {
		fmt.Println("User already exist")
		return -1
	}
	insertStatement := `insert into "users" ("username") 
	values (${1})`
	_, err = db.Exec(insertStatement, d.UserName)
	if err != nil {
		fmt.Println(err)
		return -1
	}
	userID = userExist(d.UserName)
	if userID == -1 {
		return userID
	}
	insertStatement = `insert into "userdata" ("userid", "name", "surname", "description")
	values ($1, $2, $3, $4)`
	_, err = db.Exec(insertStatement, userID, d.Name, d.Surname, d.Description)
	if err != nil {
		fmt.Println("db.Exec()", err)
		return -1
	}

	return userID
}

func DeleteUser(id int) error {
	db, err := OpenConnection([]string{})
	if err != nil {
		return err
	}
	defer db.Close()
	statement := fmt.Sprintf(`select "username" from "users" where id = %d`, id)
	rows, err := db.Query(statement)
	if err != nil {
		return err
	}
	var username string
	for rows.Next() {
		err = rows.Scan(&username)
		if err != nil {
			return err
		}
	}
	defer rows.Close()
	if userExist(username) != id {
		return fmt.Errorf("User with ID %d does not exist", id)
	}
	deleteStatement := `delete from "userdata" where userid=$1`
	_, err = db.Exec(deleteStatement, id)
	if err != nil {
		return err
	}
	deleteStatement = `delete from "users" where id=$1`
	_, err = db.Exec(deleteStatement, id)
	if err != nil {
		return err
	}
	return nil
}

func ListUsers() ([]UserData, error) {
	db, err := OpenConnection([]string{})
	if err != nil {
		return nil, err
	}
	defer db.Close()
	Data := []UserData{}
	statement := `select "id", "username", "name", "surname", "description" from "users", "userdata" where user.id = userdata.userId`
	rows, err := db.Query(statement)
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		var id int
		var username string
		var name string
		var surname string
		var description string
		err = rows.Scan(&id, &username, &name, &surname, &description)
		if err != nil {
			return Data, err
		}
		temp := UserData{ID: id, UserName: username, Name: name, Surname: surname, Description: description}
		Data = append(Data, temp)
	}
	defer rows.Close()
	return Data, nil
}
