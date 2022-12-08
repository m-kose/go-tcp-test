package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"strings"
)

type User struct {
	name     string
	conn     net.Conn
	commands chan command
}

func (user *User) readInput() {
	for {
		data, err := bufio.NewReader(user.conn).ReadString('\n')
		if err != nil {
			if strings.Contains(err.Error(), "connection closed") {
				user.conn.Close()
				return
			} else {
				log.Fatal(err)
			}
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
		case "/files":
			user.commands <- command{
				id:   CMD_FILE,
				User: user,
				args: args,
			}
		default:
			user.conn.Write([]byte("Unknown command\n"))
		}
	}
}

func ConnectToServer(a []string) {
	conn, err := net.Dial(CONN_TYPE, a[2]+":"+a[3])
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	go func() {
		for {
			data, err := bufio.NewReader(conn).ReadString('\n')
			if err != nil {
				if strings.Contains(err.Error(), "connection closed") {
					conn.Close()
					return
				} else {
					log.Fatal(err)
				}
			}

			data = strings.Trim(data, "\r\n")
			args := strings.Split(data, " ")
			command := strings.TrimSpace(args[0])
			switch command {
			case "/name":
				// Send a CMD_NAME command to the server
				conn.Write([]byte("/name " + strings.Join(args[1:], " ") + "\n"))
			case "/users":
				// Send a CMD_LIST_USERS command to the server
				conn.Write([]byte("/users\n"))
			case "/quit":
				// Send a CMD_QUIT command to the server
				conn.Write([]byte("/quit\n"))
			case "/files":
				// Send a CMD_FILE command to the server
				conn.Write([]byte("/files\n"))
			default:
				fmt.Println("Unknown command")
			}
		}
	}()
}
