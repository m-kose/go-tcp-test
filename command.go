package main

type commandID int

const (
	CMD_NAME commandID = iota
	CMD_LIST_USERS
	CMD_QUIT
)

type command struct {
	id   commandID
	User *User
	args []string
}
