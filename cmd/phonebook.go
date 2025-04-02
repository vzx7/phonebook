package main

import (
	"database/sql"
	"fmt"
	"os"
	"strings"

	"github.com/vzx7/phonebook/pkg/db"
)

func main() {
	arguments := os.Args
	if len(arguments) == 1 {
		fmt.Println("Usage: insert|delete|search|list <arguments>")
		return
	}

	if err := SetCSV(); err != nil {
		fmt.Println(err)
	}

	fileInfo, err := os.Stat(CSVFILE)

	if err != nil {
		fmt.Println(err)
		return
	}

	mode := fileInfo.Mode()
	if !mode.IsRegular() {
		fmt.Println(CSVFILE, "not regular file!")
		return
	}

	err = ReadCSVFile(CSVFILE)
	if err != nil {
		fmt.Println(err)
		return
	}

	err = CreateIndex()
	if err != nil {
		fmt.Println("Cannot create index")
	}

	// Differentiate between the commands
	switch arguments[1] {
	case "search":
		if len(arguments) != 3 {
			fmt.Println("Usage: search Number")
			return
		}
		t := strings.ReplaceAll(arguments[2], "-", "")
		if !MatchTel(t) {
			fmt.Println("Not a valid telephone number:", t)
			return
		}
		temp := Search(t)
		if temp == nil {
			fmt.Println("Number not found:", t)
			return
		}
		fmt.Println(*temp)
	case "delete":
		if len(arguments) != 3 {
			fmt.Println("Usage: delete Number")
			return
		}
		t := strings.ReplaceAll(arguments[2], "-", "")
		if !MatchTel(t) {
			fmt.Println("Not a valid telephone number:", t)
			return
		}
		err := DeleteEntry(t)
		if err != nil {
			fmt.Println(err)
		}
	case "insert":
		if len(arguments) != 5 {
			fmt.Println("Usage: insert Name Surname Telephone")
			return
		}
		t := strings.ReplaceAll(arguments[4], "-", "")
		if !MatchTel(t) {
			fmt.Println("Not a valid telephone number:", t)
			return
		}
		temp := InitS(arguments[2], arguments[3], t)
		if temp != nil {
			err = Insert(temp)
			if err != nil {
				fmt.Println(err)
				return
			}
		}
	case "list":
		List()
	default:
		fmt.Println("Not a valid option")
	}
	dbWork()
}

func dbWork() {
	dbConn, err := db.OpenConnection(db.ConnectSet{})
	if err != nil {
		fmt.Println("Ошибка подключения к БД:", err)
		return
	}
	defer dbConn.Close()
	randomUsername := RandomString(5)
	id, err := addUser(dbConn, randomUsername)
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Printf("Пользователь: %s добавлен c id: %d\n", randomUsername, id)
	}
	err = listUser(dbConn)
	if err != nil {
		fmt.Println(err)
	}
	u, err := db.GetUser(dbConn, id)
	if err != nil {
		fmt.Println(err)
	}
	u.Description = "Совершенно новый пользователь"
	err = updateUser(dbConn, u)
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Printf("Пользователь: %v обновнлен\n", u)
	}
	err = deleteUser(dbConn, u.ID)
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Printf("Пользователь c id: %d успешно удален\n", u.ID)
	}
}

func listUser(dbConn *sql.DB) error {
	data, err := db.ListUsers(dbConn)
	if err != nil {
		return fmt.Errorf("oшибка получения списка пользователей: %v", err)
	}

	for _, v := range data {
		fmt.Println(v)
	}

	return nil
}

func addUser(dbConn *sql.DB, userName string) (int, error) {
	t := db.UserData{
		UserName:    userName,
		Name:        "Vasilii",
		Surname:     "Zaz",
		Description: "This is me!",
	}

	id, err := db.AddUser(dbConn, t)
	if err != nil {
		return -1, fmt.Errorf("oшибка добавления пользователя: %s: %v", userName, err)
	}

	return id, nil
}

func deleteUser(dbConn *sql.DB, id int) error {
	err := db.DeleteUser(dbConn, id)
	if err != nil {
		return fmt.Errorf("oшибка удаления пользователя c id: %d: %v", id, err)
	}
	return nil
}

func updateUser(dbConn *sql.DB, userData *db.UserData) error {
	err := db.UpdateUser(dbConn, *userData)
	if err != nil {
		return fmt.Errorf("oшибка обновления пользователя c id: %d: %v", userData.ID, err)
	}
	return nil
}
