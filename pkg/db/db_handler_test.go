package db

import (
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
)

func TestUserExist(t *testing.T) {
	mockDB, mock, _ := sqlmock.New()
	defer mockDB.Close()

	mock.ExpectQuery(`SELECT id FROM users WHERE username = \$1`).WithArgs("johndoe").WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))
	id, _ := userExist(mockDB, "johndoe")
	if id != 1 {
		t.Errorf("Ожидался ID 1, получено %d", id)
	}
}

func TestListUsers(t *testing.T) {
	mockDB, mock, _ := sqlmock.New()
	defer mockDB.Close()

	rows := sqlmock.NewRows([]string{"id", "username", "name", "surname", "description"}).AddRow(1, "johndoe", "John", "Doe", "Sample user")
	mock.ExpectQuery("SELECT u.id, u.username, ud.name, ud.surname, ud.description FROM users u").WillReturnRows(rows)
	users, _ := ListUsers(mockDB)
	if len(users) != 1 || users[0].ID != 1 {
		t.Errorf("Ошибка получения списка пользователей")
	}
}

func TestUpdateUser(t *testing.T) {
	mockDB, mock, _ := sqlmock.New()
	defer mockDB.Close()

	mock.ExpectExec(`UPDATE userdata SET name=\$1, surname=\$2, description=\$3 WHERE userid=\$4`).WithArgs("John", "Doe", "Updated description", 1).WillReturnResult(sqlmock.NewResult(1, 1))
	err := UpdateUser(mockDB, UserData{ID: 1, Name: "John", Surname: "Doe", Description: "Updated description"})
	if err != nil {
		t.Errorf("Ошибка при обновлении пользователя: %v", err)
	}
}

func TestDeleteUser(t *testing.T) {
	mockDB, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Error initializing sqlmock: %s", err)
	}
	defer mockDB.Close()

	mock.ExpectExec(`DELETE FROM userdata WHERE userid=\$1`).WithArgs(1).WillReturnResult(sqlmock.NewResult(0, 1))
	mock.ExpectExec(`DELETE FROM users WHERE id=\$1`).WithArgs(1).WillReturnResult(sqlmock.NewResult(0, 1))

	err = DeleteUser(mockDB, 1)
	if err != nil {
		t.Errorf("Unexpected error: %s", err)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("Unfulfilled expectations: %s", err)
	}
}

func TestGetUser(t *testing.T) {
	mockDB, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Error initializing sqlmock: %s", err)
	}
	defer mockDB.Close()

	mock.ExpectQuery(`SELECT u.id, u.username, ud.name, ud.surname, ud.description FROM users u JOIN userdata ud ON u.id = ud.userid WHERE u.id = \$1`).WithArgs(1).WillReturnRows(sqlmock.NewRows([]string{"id", "username", "name", "surname", "description"}).AddRow(1, "johndoe", "John", "Doe", "Test user"))

	user, err := GetUser(mockDB, 1)
	if err != nil {
		t.Errorf("Unexpected error: %s", err)
	}
	if user == nil {
		t.Errorf("Expected user data, got nil")
	}
	if user.ID != 1 || user.UserName != "johndoe" || user.Name != "John" || user.Surname != "Doe" || user.Description != "Test user" {
		t.Errorf("User data does not match expected values")
	}
}
