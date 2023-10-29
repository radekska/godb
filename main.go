package main

import (
	"fmt"
	godb "godb/pkg"
	"os"
	"strconv"
)

func main() {
	args := os.Args
	if len(args) != 6 {
		fmt.Println("Usage: main.go <host> <port> <username> <password> <database>")
		return
	}
	godb.Host = args[1]
	port, err := strconv.Atoi(args[2])
	if err != nil {
		fmt.Println("Error: port must be an integer")
		return
	}
	godb.Port = port
	godb.Username = args[3]
	godb.Password = args[4]
	godb.Database = args[5]

	db, err := godb.OpenConnection()

	userID, err := godb.AddUser(db, godb.User{Username: "test", Name: "Test", Surname: "Test", Description: "Test"})
	if err != nil {
		fmt.Println("Error adding user:", err)
		return
	}
	fmt.Println(godb.ListUsers(db))

	err = godb.UpdateUser(
		db,
		godb.User{ID: userID, Username: "test", Name: "Test2", Surname: "Test2", Description: "Test2"},
	)
	if err != nil {
		fmt.Println("Error updating user:", err)
		return
	}
	fmt.Println(godb.ListUsers(db))

	err = godb.DeleteUser(db, userID)
	if err != nil {
		fmt.Println("Error deleting user:", err)
		return
	}
	fmt.Println(godb.ListUsers(db))
}
