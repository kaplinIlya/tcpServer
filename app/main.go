package main

import (
	rotatelogs "github.com/lestrrat-go/file-rotatelogs"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"io"
	"net"
	"os"
	"tcpServer/handlers"
	"time"
)

func main() {
	if err := initConfig(); err != nil {
		logrus.Fatalf("Config initialization error: %s", err.Error())
	}

	logFilePattern := viper.GetString("app.log.filenamePattern")
	if logFilePattern == "" {
		panic("empty app.log.filenamePattern in config")
	}
	writer, _ := rotatelogs.New(
		logFilePattern+".%Y-%m-%d",
		//logPath+".%H%M%S",
		rotatelogs.WithLinkName(logFilePattern),
		rotatelogs.WithMaxAge(time.Duration(viper.GetInt32("app.log.maxAgeInHours"))*time.Hour),
		rotatelogs.WithRotationTime(time.Duration(viper.GetInt32("app.log.rotationInHours"))*time.Hour),
	)
	logrus.SetFormatter(new(logrus.JSONFormatter))
	logrus.SetOutput(io.MultiWriter(writer, os.Stdout))
	defer writer.Close()

	port := ":" + viper.GetString("app.port")
	listener, err := net.Listen("tcp", port)
	if err != nil {
		panic(err)
	}
	logrus.Infof("starting server on port %s", port)
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
