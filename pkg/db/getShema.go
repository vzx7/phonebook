package db

import (
	"database/sql"
	"fmt"
	"strconv"

	_ "github.com/lib/pq"
)

func Connect(arguments []string) {
	/* 	arguments := os.Args */
	if len(arguments) != 6 {
		fmt.Println("Please provide: hostname port username password db!")
		return
	}

	_host := arguments[1]
	_port := arguments[2]
	_user := arguments[3]
	_password := arguments[4]
	_dbName := arguments[5]

	port, err := strconv.Atoi(_port)
	if err != nil {
		fmt.Println("Not a valid port number:", err)
		return
	}

	conn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", _host, port, _user, _password, _dbName)
	db, err := sql.Open("postgres", conn)
	if err != nil {
		fmt.Println("Open():", err)
	}
	//fmt.Println(db.Stats().InUse)
	defer db.Close()

	rows, err := db.Query(`SELECT "datname" FROM "pg_database"
	WHERE datistemplate = false`)
	if err != nil {
		fmt.Println("Query 1:", err)
		return
	}

	for rows.Next() {
		var name string
		err = rows.Scan(&name)
		if err != nil {
			fmt.Println("Scan:", err)
			return
		}
		fmt.Println("*", name)
	}
	defer rows.Close()

	rows, err = db.Query(`SELECT table_name FROM information_schema.tables WHERE 
	table_schema = 'public' ORDER BY table_name`)
	if err != nil {
		fmt.Println("Query", err)
		return
	}

	for rows.Next() {
		var name string
		err = rows.Scan(&name)
		if err != nil {
			fmt.Println("Scan", err)
			return
		}
		fmt.Println("+T", name)
	}
	defer rows.Close()
}
