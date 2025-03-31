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
	statment := fmt.Sprintf(`SELECT "id" FROM "users" WHERE username = %s`, username)
	rows, err := db.Query(statment)
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
	statment := fmt.Sprintf(`select "username" from "users" where id = %d`, id)
	rows, err := db.Query(statment)
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
	deletestatment := `delete from "userdata" where userid=$1`
	_, err = db.Exec(deletestatment, id)
	if err != nil {
		return err
	}
	deletestatment = `delete from "users" where id=$1`
	_, err = db.Exec(deletestatment, id)
	if err != nil {
		return err
	}
	return nil
}
