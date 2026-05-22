package constructor

import (
	"answer_protocol/internal/models"
	"net"
)

func NewHub() *models.Hub{
	return &models.Hub{
		Register:   make(chan net.Conn),
		Unregister: make(chan net.Conn), 
		Broadcast:  make(chan string),
		Clients: 	make(map[net.Conn]*models.Player),
	}
}

func NewPlayer(conn string, name string) *models.Player{
	return &models.Player{
		Id: 	conn,
		Name: 	name,
	}
}
