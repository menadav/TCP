package network

import (
	"answer_protocol/src/constructor"
	"answer_protocol/src/logger"
	"answer_protocol/src/models"
	"answer_protocol/src/parse"
	"bufio"
	"net"
	"time"
)

const floodThreshold = 20

func ClientAtender(conn net.Conn, hub *models.Hub) {
	scanner := bufio.NewScanner(conn)
	name_p := ""
	name_p = Authentication(scanner, conn, hub)
	if name_p == "" {
		logger.Warn("auth failed", "addr", logger.Addr(conn))
		conn.Close()
		return
	}
	start := hub.World.Rooms["start"]
	player := constructor.NewPlayer(conn.RemoteAddr().String(), conn, name_p, start)
	if start != nil {
		start.Mu.Lock()
		if start.Players == nil {
			start.Players = make(map[string]*models.Player)
		}
		start.Players[player.Name] = player
		start.Mu.Unlock()
	}
	defer func() {
		if name_p != "" {
			if player.Room != nil {
				player.Room.Mu.Lock()
				delete(player.Room.Players, player.Name)
				player.Room.Mu.Unlock()
			}
			hub.Unregister <- player
		}
		conn.Close()
		logger.Info("connection close", "name", name_p, "addr", logger.Addr(conn))
	}()
	hub.Register <- player
	logger.Info("client registered", "name", name_p, "addr", logger.Addr(conn))
	go player.ListenMsg()
	cmdCount := 0
	windowStart := time.Now()
	processFunction := func(line string) {
		cmdCount++
		if time.Since(windowStart) > 10*time.Second {
			cmdCount = 1
			windowStart = time.Now()
		} else if cmdCount == floodThreshold {
			logger.Warn("abuse detected", "name", name_p, "addr", logger.Addr(conn), "reason", "command_flood", "count", cmdCount, "window", "10s")
		}
		parse.ParseCommandCli(line, player, hub)
	}
	StartScanner(scanner, processFunction, conn)
	logger.Info("client pick control + D", "name", name_p, "addr", logger.Addr(conn))
}
