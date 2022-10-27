package main

import (
	"fmt"
	"log"
	"net"
	"os"
	"time"
)

//https://github.com/pplam/tcp-file-transfer/blob/master/server/main.go
//https://www.youtube.com/watch?v=Sphme0BqJiY

const (
	CONN_HOST = "localhost"
	CONN_TYPE = "tcp"
)

type server struct {
	commands chan command
}

func newServer() *server {
	return &server{
		commands: make(chan command),
	}
}

func (s *server) run() {
	for cmd := range s.commands {
		switch cmd.id {
		case CMD_NAME:
			s.name0(cmd.User, cmd.args)
		case CMD_LIST_USERS:
			s.listall0(cmd.User, cmd.args)
		case CMD_QUIT:
			s.quit0(cmd.User, cmd.args)
		}
	}
}

func (s *server) name0(u *User, args []string) {
	u.name = args[1]
	u.conn.Write([]byte(fmt.Sprintf("Your name is set to be: %s", u.name)))
}

func (s *server) listall0(u *User, args []string) {

}

func (s *server) quit0(u *User, args []string) {

}

func (s *server) newUser(conn net.Conn) {
	log.Printf("New user has connected: %s", conn.RemoteAddr().String())
	u := &User{
		name:     "anonymous",
		conn:     conn,
		commands: s.commands,
	}
	u.readInput()
}

func CreateServer(a []string) {
	s := newServer()
	go s.run()
	listener, err := net.Listen(CONN_TYPE, CONN_HOST+":"+a[2])
	if err != nil {
		log.Fatal(err)
	}
	defer listener.Close()
	for {
		conn, _ := listener.Accept()
		go s.newUser(conn)
		for {
			receiveFile(conn)
			fmt.Println("Checking for new nodes...")
			time.Sleep(15 * time.Second)
			continue
		}
	}

	//sendFile(server)
}

func receiveFile(conn net.Conn) {
	//conn, _ := server.Accept()
	//defer conn.Close()
	path, _ := os.Getwd()
	if err := os.Mkdir(path+"\\"+"a", 0777); err != nil && !os.IsExist(err) {
		log.Fatal(err)
	}
}
