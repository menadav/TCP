package main

import (
	"answer_protocol/src/constructor"
	"answer_protocol/src/logger"
	"answer_protocol/src/network"
	"answer_protocol/src/speakserver"
	"answer_protocol/src/world"
	"net"
)

func main() {
	listen, err := net.Listen("tcp", ":8080")
	if err != nil {
		logger.Error("listen failed", "addr", ":8080", "err", err)
		return
	}
	defer listen.Close()
	data, err := world.LoadWorld("data.yaml")
	if err != nil {
		logger.Error("world load failed", "path", "data.yaml", "err", err)
		return
	}
	hub := constructor.NewHub(data)
	go hub.Run()
	logger.Info("server ready", "addr", ":8080")
	for {
		conn, err := listen.Accept()
		if err != nil {
			logger.Error("accept failed", "err", err)
			continue
		}
		logger.Info("connection open", "addr", logger.Addr(conn))
		go network.ClientAtender(conn, hub)
		speak.SendSuccess(conn, "hello proto=1")
	}
}
