package main

import (
	"fmt"
	"log"
	"net"
	"os"
	"strings"
	"time"
)

//https://github.com/pplam/tcp-file-transfer/blob/master/server/main.go
//https://www.youtube.com/watch?v=Sphme0BqJiY

const (
	CONN_HOST = "localhost"
	CONN_TYPE = "tcp"
)

type server struct {
	commands map[net.Conn]chan command
	users    []*User
}

func newServer() *server {
	return &server{
		commands: make(map[net.Conn]chan command),
	}
}

func (s *server) name0(u *User, args []string) string {
	u.name = args[1]
	u.conn.Write([]byte(fmt.Sprintf("Your name is set to be: %s\n", u.name)))
	return args[1]
}

// Write a listall0 function that lists all the users' names.

func (s *server) quit0(u *User, args []string) {
	u.conn.Write([]byte("Bye!\n"))
	for i, user := range s.users {
		if user == u {
			s.users = append(s.users[:i], s.users[i+1:]...)
			break
		}
	}
}

var users [100]string
var slice = users[0:len(users)]

// Write a newUser function that creates a new user and adds it to the server.
func (s *server) newUser(conn net.Conn) {
	addr := conn.RemoteAddr().String()
	name := strings.Split(addr, ":")[0]
	u := &User{
		name:     name,
		conn:     conn,
		commands: make(chan command),
	}
	s.commands[conn] = u.commands
	s.users = append(s.users, u)

	go u.readInput()
	go s.processCommands(u)
}

func (s *server) listall0(u *User, users []string) {
	u.conn.Write([]byte("List of users:\n"))
	for _, user := range s.users {
		name := strings.TrimSpace(user.name)
		u.conn.Write([]byte(name + "\n"))
	}
}

func (s *server) processCommands(u *User) {
	conn := u.conn
	for {
		select {
		case cmd, ok := <-s.commands[conn]:
			// Check if the channel has been closed
			if !ok {
				// If the channel has been closed, remove the user from the server
				delete(s.commands, cmd.User.conn)
				continue
			}

			// Process the command
			switch cmd.id {
			case CMD_NAME:
				s.name0(cmd.User, cmd.args)
			case CMD_LIST_USERS:
				s.listall0(cmd.User, cmd.args)
			case CMD_QUIT:
				s.quit0(cmd.User, cmd.args)
			default:
				cmd.User.conn.Write([]byte("Unknown command\n"))
			}
		}
	}
}

func CreateServer(a []string) {
	s := newServer()
	listener, err := net.Listen(CONN_TYPE, CONN_HOST+":"+a[2])
	slice = append(slice, listener.Addr().String())
	if err != nil {
		log.Fatal(err)
	}
	defer listener.Close()
	for {
		conn, _ := listener.Accept()
		go s.newUser(conn)
		for {
			fmt.Println("Checking for new nodes...\n")
			time.Sleep(15 * time.Second)
			continue
		}
	}
}

// Write a receiveFile function that receives tx.db from the server.
func (s *server) receiveFile(conn net.Conn, args []string) {
	// Create a buffer to store the file in
	buf := make([]byte, 1024)
	// Open a file for writing
	f, err := os.Create("tx.db")
	if err != nil {
		fmt.Println(err)
		return
	}
	defer f.Close()
	// Receive the file
	for {
		_, err := conn.Read(buf)
		if err != nil {
			fmt.Println(err)
			return
		}
		// Write the received bytes to the file
		_, err = f.Write(buf)
		if err != nil {
			fmt.Println(err)
			return
		}
	}
}
