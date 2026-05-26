package constructor

import (
	"answer_protocol/src/models"
	"net"
)

func NewHub() *models.Hub{
	return &models.Hub{
		Register:   make(chan *models.Player),
		Unregister: make(chan *models.Player),
		Broadcast:  make(chan models.Message),
		Clients: 	make(map[net.Conn]*models.Player),
	}
}

func NewPlayer(conn_st string, conn net.Conn, name string, startRoom string) *models.Player{
	return &models.Player{
		Id: 	conn_st,
		Conn:	conn,
		Name: 	name,
		Room:	startRoom,
		Group:	"",
	}
}
