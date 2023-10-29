package godb

import (
	"database/sql"
	"errors"
	"fmt"

	_ "github.com/lib/pq"
)

var (
	Host     = "localhost"
	Port     = 5432
	Username = "postgres"
	Password = "postgres"
	Database = "postgres"
)

const (
	usersTable     = "users"
	usersDataTable = "userdata"
)

type User struct {
	ID          string
	Username    string
	Name        string
	Surname     string
	Description string
}

func OpenConnection() (*sql.DB, error) {
	connStr := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		Host, Port, Username, Password, Database)
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, err
	}
	return db, nil
}

func userExists(db *sql.DB, username string) (string, error) {
	var id string
	err := db.QueryRow(fmt.Sprintf("SELECT id FROM %s WHERE username = $1", usersTable), username).Scan(&id)
	if err != nil {
		return "", err
	}
	return id, nil
}

func DeleteUser(db *sql.DB, userID string) error {
	_, err := db.Exec(fmt.Sprintf("DELETE FROM %s WHERE id = $1", usersTable), userID)
	if err != nil {
		return err
	}
	_, err = db.Exec(fmt.Sprintf("DELETE FROM %s WHERE userid = $1", usersDataTable), userID)
	return err
}

func UpdateUser(db *sql.DB, user User) error {
	userID, err := userExists(db, user.Username)
	if err != nil {
		return err
	}
	if userID == "" {
		return errors.New("User does not exist")
	}
	_, err = db.Exec(
		fmt.Sprintf("UPDATE %s SET name = $1, surname = $2, description = $3 WHERE userid = $4", usersDataTable),
		user.Name,
		user.Surname,
		user.Description,
		userID,
	)
	return err
}

func AddUser(db *sql.DB, user User) (string, error) {
	_, err := db.Exec(fmt.Sprintf("INSERT INTO %s (username) VALUES ($1)", usersTable), user.Username)
	if err != nil {
		return "", err
	}

	userID, err := userExists(db, user.Username)
	if err != nil {
		return "", err
	}
	_, err = db.Exec(
		fmt.Sprintf("INSERT INTO %s (userid, name, surname, description) VALUES ($1, $2, $3, $4)",
			usersDataTable),
		userID,
		user.Name,
		user.Surname,
		user.Description,
	)
	if err != nil {
		return "", err
	}
	return userID, nil
}

func ListUsers(db *sql.DB) ([]User, error) {
	rows, err := db.Query(
		fmt.Sprintf(
			"SELECT id, username, name, surname, description FROM %s INNER JOIN %s ON %s.id = %s.userid",
			usersTable,
			usersDataTable,
			usersTable,
			usersDataTable,
		),
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	users := make([]User, 0)
	for rows.Next() {
		var user User
		err = rows.Scan(&user.ID, &user.Username, &user.Name, &user.Surname, &user.Description)
		if err != nil {
			return nil, err
		}
		users = append(users, user)
	}
	return users, nil
}
