package main

import (
	"bufio"
	"log"
	"net"
	"strings"
)

type User struct {
	name     string
	conn     net.Conn
	commands chan<- command
}

func (user *User) readInput() {
	for {
		data, err := bufio.NewReader(user.conn).ReadString('\n')
		if err != nil {
			log.Fatal(err)
		}

		data = strings.Trim(data, "\r\n")
		args := strings.Split(data, " ")
		Command := strings.TrimSpace(args[0])
		switch Command {
		case "/name":
			user.commands <- command{
				id:   CMD_NAME,
				User: user,
				args: args,
			}
		case "/users":
			user.commands <- command{
				id:   CMD_LIST_USERS,
				User: user,
				args: args,
			}
		case "/quit":
			user.commands <- command{
				id:   CMD_QUIT,
				User: user,
				args: args,
			}
		default:
			user.conn.Write([]byte("Unknown command"))
		}
	}
}

func ConnectToServer(a []string) {
	_, err := net.Dial(CONN_TYPE, a[2]+":"+a[3])
	if err != nil {
		log.Fatal(err)
	}
}
