package game

import (
    "answer_protocol/src/models"
    "answer_protocol/src/speakserver"
    "strings"
)

func TakeItem(player *models.Player, query string) {
    actualRoom := player.Room
    actualRoom.Mu.Lock()
    defer actualRoom.Mu.Unlock()

    for i, item := range actualRoom.Items {
        if strings.EqualFold(item.Name, query) || item.ID == query {
            if !item.Obtainable {
                speak.SendError(player.Conn, 403, "ITEM_NOT_OBTAINABLE")
                return
            }
            actualRoom.Items = append(actualRoom.Items[:i], actualRoom.Items[i+1:]...)
            player.Inventory = append(player.Inventory, item)
            speak.SendSuccess(player.Conn, "taken="+item.ID)
            return
        }
    }
    speak.SendError(player.Conn, 404, "ITEM_NOT_FOUND")
}