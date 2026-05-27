package utils

import (
	"net"
	"answer_protocol/src/models"

)

func ExistName(clients map[net.Conn]*models.Player, name string) bool {
	for _, player := range clients {
		if player.Name == name {
			return true
		}
	}
	return false
}
