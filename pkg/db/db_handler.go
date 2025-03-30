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
