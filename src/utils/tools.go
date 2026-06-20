package utils

import (
	"answer_protocol/src/models"
	"net"
)

func ExistName(clients map[net.Conn]*models.Player, name string) bool {
	for _, player := range clients {
		if player.Name == name {
			return true
		}
	}
	return false
}
