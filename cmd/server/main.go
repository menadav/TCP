package main

import (
	"answer_protocol/src/constructor"
	"answer_protocol/src/logger"
	"answer_protocol/src/network"
	"answer_protocol/src/speakserver"
	"answer_protocol/src/world"
	"net"
	"time"
)

const (
	rapidConnThreshold = 5
	rapidConnWindow    = 10 * time.Second
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
	connSeen := make(map[string][]time.Time)
	for {
		conn, err := listen.Accept()
		if err != nil {
			logger.Error("accept failed", "err", err)
			continue
		}
		addr := logger.Addr(conn)
		logger.Info("connection open", "addr", addr)
		if host, _, e := net.SplitHostPort(addr); e == nil {
			now := time.Now()
			recent := connSeen[host][:0]
			for _, t := range connSeen[host] {
				if now.Sub(t) <= rapidConnWindow {
					recent = append(recent, t)
				}
			}
			recent = append(recent, now)
			connSeen[host] = recent
			if len(recent) > rapidConnThreshold {
				logger.Warn("abuse detected", "addr", addr, "reason", "rapid_connection", "count", len(recent), "window", "10s")
			}
		}
		go network.ClientAtender(conn, hub)
		speak.SendSuccess(conn, "hello proto=1")
	}
}
