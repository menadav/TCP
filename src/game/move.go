package game

import (
	"answer_protocol/src/models"
	"answer_protocol/src/speakserver"
)

func MapMove(player *models.Player, move string, hub *models.Hub) {
	oldRoomID := player.Room.Id
	err := hub.World.MovePlayer(player, move)
	if err != nil {
		speak.SendErr(player.Conn, speak.ErrNoExit)
		return
	}
	player.HandleRoomVisit(player.Room.Id, hub.World.Quest)
	hub.NotifyRoomLeave(player, oldRoomID)
	hub.NotifyRoomEnter(player, player.Room.Id)
	speak.SendSuccess(player.Conn, "room="+player.Room.Id)
}
