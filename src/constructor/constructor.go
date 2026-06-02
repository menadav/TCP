package constructor

import (
	"answer_protocol/src/models"
	"net"
)

func NewWorld() *models.World{
	return &models.World{
		Rooms:	make(map[string]*models.Room),
		Items:	make(map[string]*models.Item),
		Npcs:	make(map[string]*models.Npc),
	}
}

func NewHub(data *models.World) *models.Hub{
	return &models.Hub{
		Register:   make(chan *models.Player),
		Unregister: make(chan *models.Player),
		Broadcast:  make(chan models.Message),
		Clients: 	make(map[net.Conn]*models.Player),
		Groups:     make(map[string]*models.Group),
		World:		data,
	}
}

func NewPlayer(conn_st string, conn net.Conn, name string, startRoom *models.Room) *models.Player{
	return &models.Player{
		Id: 	conn_st,
		Conn:	conn,
		Name: 	name,
		Room:	startRoom,
		Group:	"",
		Inventory: []*models.Item{},
	}
}
