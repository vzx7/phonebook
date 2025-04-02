package db

import (
	"database/sql"
	"fmt"
)

/*
OpenConnection - функция для подключения к БД postgres.
Если не передан аргумент []string, с данными для подключения,
то будет использован данные с pkg/db/constants.go
*/
func OpenConnection(connectSet ConnectSet) (*sql.DB, error) {
	_host := CONNECT_SET.Host
	_port := CONNECT_SET.Port
	_user := CONNECT_SET.DBUser
	_password := CONNECT_SET.DBPassword
	_dbName := CONNECT_SET.DBName

	if connectSet.Host == "" {
		fmt.Println("The host is not transferred to the data connection function, constants will be used to connect to the database ...")
		_host = CONNECT_SET.Host
	}

	if connectSet.Port == 0 {
		fmt.Println("The port is not transferred to the data connection function, constants will be used to connect to the database ...")
		_port = CONNECT_SET.Port
	}
	if connectSet.DBName == "" {
		fmt.Println("The BDName is not transferred to the data connection function, constants will be used to connect to the database ...")
		_dbName = CONNECT_SET.DBName
	}
	if connectSet.DBUser == "" {
		fmt.Println("The DBUser is not transferred to the data connection function, constants will be used to connect to the database ...")
		_user = CONNECT_SET.DBUser
	}
	if connectSet.DBPassword == "" {
		fmt.Println("The DBPassword is not transferred to the data connection function, constants will be used to connect to the database ...")
		_password = CONNECT_SET.DBPassword
	}

	conn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", _host, _port, _user, _password, _dbName)
	db, err := sql.Open("postgres", conn)
	if err != nil {
		return nil, err
	}

	return db, nil
}
