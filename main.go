package main

import (
	"database/sql"
	"fmt"
	"os"
	"strconv"

	_ "github.com/lib/pq"
)

const (
	usersTable     = "users"
	usersDataTable = "userdata"
)

type User struct {
	ID          int
	Username    string
	Name        string
	Surname     string
	Description string
}
// TODO: add update, delete and fix list all

func userExists(db *sql.DB, username string) (string, error) {
	var id string
	err := db.QueryRow(fmt.Sprintf("SELECT id FROM %s WHERE username = $1", usersTable), username).Scan(&id)
	if err != nil {
		return "", err
	}
	return id, nil
}

func AddUser(db *sql.DB, user User) error {
	_, err := db.Exec(fmt.Sprintf("INSERT INTO %s (username) VALUES ($1)", usersTable), user.Username)
	if err != nil {
		return err
	}

	userID, err := userExists(db, user.Username)
	if err != nil {
		return err
	}
	_, err = db.Exec(
		fmt.Sprintf("INSERT INTO %s (userid, name, surname, description) VALUES ($1, $2, $3, $4)",
			usersDataTable),
		userID,
		user.Name,
		user.Surname,
		user.Description,
	)
	return err
}

func getUsers(db *sql.DB) ([]User, error) {
	rows, err := db.Query(fmt.Sprintf("SELECT * FROM %s", usersTable))
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var users []User
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

func main() {
	args := os.Args
	if len(args) != 6 {
		fmt.Println("Usage: go run main.go <host> <port> <user> <password> <database>")
		return
	}
	host, portS, user, password, database := args[1], args[2], args[3], args[4], args[5]
	port, err := strconv.Atoi(portS)
	if err != nil {
		fmt.Println("Port must be a number", err)
		return
	}
	conn := fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		host,
		port,
		user,
		password,
		database,
	)
	db, err := sql.Open("postgres", conn)
	if err != nil {
		fmt.Println("Error connecting to database", err)
		return
	}
	defer db.Close()

	rows, err := db.Query(`SELECT "datname" FROM "pg_database"`)
	if err != nil {
		fmt.Println("Error querying database", err)
		return
	}
	for rows.Next() {
		var name string
		err = rows.Scan(&name)
		if err != nil {
			fmt.Println("Error scanning row", err)
			return
		}
		fmt.Println(name)
	}

	query := `SELECT table_name FROM information_schema.tables WHERE table_schema = 'public' ORDER BY table_name`
	rows, err = db.Query(query)
	if err != nil {
		fmt.Println("Error querying database", err)
		return
	}
	for rows.Next() {
		var name string
		err = rows.Scan(&name)
		if err != nil {
			fmt.Println("Error scanning row", err)
			return
		}
		fmt.Println("+T", name)
	}
	defer rows.Close()

	err = AddUser(db, User{Username: "john123", Name: "John", Surname: "Doe", Description: "A user"})
	if err != nil {
		fmt.Println("Error adding user", err)
		return
	}

	// users, err := getUsers(db)
	// if err != nil {
	// 	fmt.Println("Error getting users", err)
	// 	return
	// }
	// fmt.Println(users)
}
