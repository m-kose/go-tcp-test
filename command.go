package main

type commandID int

const (
	CMD_NAME commandID = iota
	CMD_LIST_USERS
	CMD_QUIT
	CMD_NEW_USER
	CMD_FILE
)

type command struct {
	id   commandID
	User *User
	args []string
}
