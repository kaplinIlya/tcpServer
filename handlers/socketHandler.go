package handlers

import (
	"bufio"
	"fmt"
	"github.com/sirupsen/logrus"
	"net"
)

var greeting string
var cons []net.Conn

func Init(s string) {
	greeting = s
}

func AddConn(con net.Conn) {
	cons = append(cons, con)
}

func HandleConnection(conn net.Conn) {
	name := conn.RemoteAddr().String()
	logrus.Infof("%+v connected\n", name)
	conn.Write([]byte(greeting + name + "\n\r"))
	defer conn.Close()
	scanner := bufio.NewScanner(conn)
	for scanner.Scan() {
		text := scanner.Text()
		if text == "Exit" {
			conn.Write([]byte("Bye\n\r"))
			logrus.Info(fmt.Println(name, "disconnected"))
			break
		} else if text != "" {
			logrus.Infof("%s send: %s", name, text)
			for _, c := range cons {
				if c.RemoteAddr().String() == conn.RemoteAddr().String() {
					//conn.Write([]byte("You enter " + text + "\n\r"))
				} else {
					c.Write([]byte(name + " enter " + text + "\n\r"))
				}
			}
		}
	}
}
