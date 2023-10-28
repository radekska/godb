package main

import (
	"database/sql"
	"fmt"
	"os"
	"strconv"

	_ "github.com/lib/pq"
)

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
}
