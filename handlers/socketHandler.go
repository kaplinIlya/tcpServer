package handlers

import (
	"bufio"
	"encoding/json"
	"fmt"
	"github.com/sirupsen/logrus"
	"net"
)

var cons []net.Conn

type status struct {
	Count   int      `json:"count"`
	Clients []string `json:"clients"`
}

func AddConn(con net.Conn) {
	cons = append(cons, con)
}

func HandleConnection(conn net.Conn) {
	name := conn.RemoteAddr().String()
	logrus.Infof("%+v connected", name)
	conn.Write([]byte("Hello \n" +
		"exit - disconnect\n" +
		"status - get all connected clients\n" +
		"asm_upd - update asm cache notification)\n"))
	defer conn.Close()
	scanner := bufio.NewScanner(conn)
Loop:
	for scanner.Scan() {
		switch text := scanner.Text(); text {
		case "exit":
			conn.Write([]byte("Bye\r"))
			logrus.Info(fmt.Println(name, "disconnected"))
			break Loop
		case "status":
			logrus.Infof("%s request status", name)
			status := status{
				0,
				[]string{},
			}
			for _, c := range cons {
				if c.RemoteAddr().String() != conn.RemoteAddr().String() {
					status.Clients = append(status.Clients, c.RemoteAddr().String())
				}
			}
			status.Count = len(status.Clients)
			jsonRs, _ := json.Marshal(status)
			conn.Write(append(jsonRs, []byte("\n\r")...))
			break
		case "asm_upd":
			logrus.Infof("%s update asm cache", name)
			for _, c := range cons {
				if c.RemoteAddr().String() != conn.RemoteAddr().String() {
					c.Write([]byte("asm_upd\n\r"))
				}
			}
			break
		default:
			logrus.Infof("%s send: %s", name, text)
		}
	}
}
