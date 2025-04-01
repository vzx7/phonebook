package db

import (
	"database/sql"
	"fmt"
	"strconv"
)

// Open DB Connect
func OpenConnection(arguments []string) (*sql.DB, error) {
	_host := HOST
	_port := PORT
	_user := BD_USER
	_password := BD_PASSWORD
	_dbName := BD_NAME
	if len(arguments) != 5 {
		fmt.Println("The arguments are not transferred to the data connection function, constants will be used to connect to the database ...")
	} else {
		_host = arguments[0]

		port, err := strconv.Atoi(arguments[1])
		if err != nil {
			return nil, err
		}

		_port = port
		_user = arguments[2]
		_password = arguments[3]
		_dbName = arguments[4]
	}

	conn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", _host, _port, _user, _password, _dbName)
	db, err := sql.Open("postgres", conn)
	if err != nil {
		return nil, err
	}

	return db, nil
}
