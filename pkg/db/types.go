package db

type UserData struct {
	ID          int
	Name        string
	UserName    string
	Surname     string
	Description string
}

// Данные для подключения к DB
type ConnectSet struct {
	Host       string
	Port       int16
	DBUser     string
	DBPassword string
	DBName     string
}
