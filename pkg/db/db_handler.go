package db

import (
	"fmt"
	"strings"

	_ "github.com/lib/pq"
)

func exist(_username string) int {
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
