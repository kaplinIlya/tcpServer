package main

import (
	rotatelogs "github.com/lestrrat-go/file-rotatelogs"
	logrus "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"io"
	"net"
	"os"
	"tcpServer/handlers"
	"time"
)

func main() {
	const (
		port    = ":8088"
		logPath = "logs"
	)
	writer, _ := rotatelogs.New(
		logPath+".%Y-%m-%d",
		//logPath+".%H%M%S",
		rotatelogs.WithLinkName(logPath),
		rotatelogs.WithMaxAge(time.Duration(72)*time.Hour),
		rotatelogs.WithRotationTime(time.Duration(24)*time.Hour),
	)
	logrus.SetFormatter(new(logrus.JSONFormatter))
	logrus.SetOutput(io.MultiWriter(writer, os.Stdout))
	defer writer.Close()

	if err := initConfig(); err != nil {
		logrus.Fatalf("Config initialization error: %s", err.Error())
	}

	listener, err := net.Listen("tcp", port)
	logrus.Infof("starting server on port %s", port)
	if err != nil {
		panic(err)
	}
	for {
		conn, err := listener.Accept()
		if err != nil {
			panic(err)
		}
		handlers.AddConn(conn)
		go handlers.HandleConnection(conn)
	}
}

func initConfig() error {
	viper.AddConfigPath("config")
	viper.SetConfigName("config")
	return viper.ReadInConfig()
}
