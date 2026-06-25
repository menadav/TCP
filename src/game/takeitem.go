package game

import (
	"answer_protocol/src/logger"
	"answer_protocol/src/models"
	"answer_protocol/src/speakserver"
	"strings"
)

func TakeItem(player *models.Player, query string, hub *models.Hub) {
	actualRoom := player.Room
	actualRoom.Mu.Lock()

	for i, item := range actualRoom.Items {
		if strings.EqualFold(item.Name, query) || item.ID == query {
			if !item.Obtainable {
				actualRoom.Mu.Unlock()
				speak.SendErr(player.Conn, speak.ErrNotObtainable)
				return
			}
			if item.Hand && !player.Hand {
				actualRoom.Mu.Unlock()
				speak.SendErr(player.Conn, speak.ErrHandsFull)
				return
			}
			if item.Hand {
				player.UpdateDmg(item)
			}
			actualRoom.Items = append(actualRoom.Items[:i], actualRoom.Items[i+1:]...)
			player.Inventory = append(player.Inventory, item)
			player.HandleItemCollection(item.ID, hub.World.Quest)
			actualRoom.Mu.Unlock()
			speak.SendSuccess(player.Conn, "taken="+item.ID)
			logger.Info("world change", "event", "item_take", "player", player.Name, "item", item.ID, "room", actualRoom.Id)
			hub.Broadcast <- models.Message{Scope: models.ScopeRoom, Filter: actualRoom.Id, Category: "ROOM", Content: "ITEMS_CHANGED"}
			return
		}
	}
	actualRoom.Mu.Unlock()
	speak.SendErr(player.Conn, speak.ErrItemNotFound)
}
