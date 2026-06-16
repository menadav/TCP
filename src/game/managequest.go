package game

import (
    "answer_protocol/src/models"
    "answer_protocol/src/speakserver"
    "strings"
)

func ManageQuest(player *models.Player, h *models.Hub, action string, questID string){
    action = strings.ToUpper(action)

    quest, exists := h.World.Quest[questID]
    if !exists {
        speak.SendError(player.Conn, 404, "QUEST_NOT_FOUND")
        return
    }
    switch action {
    case "ACCEPT":
        var startItem *models.Item
        if quest.StartItem != ""{
            startItem = h.World.Items[quest.StartItem]
        }
        if err := player.AcceptQuest(quest, startItem); err != nil{
            speak.SendError(player.Conn, 403, err.Error())
            return
        }
        speak.SendSuccess(player.Conn, "quest_accepted="+quest.ID)
    case "COMPLETE":
        rewardItem := h.World.Items[quest.Reward]
        if err := player.CompleteQuest(quest, rewardItem); err != nil {
            speak.SendError(player.Conn, 403, err.Error())
            return
        }
        speak.SendSuccess(player.Conn, "quest_completed="+quest.ID+" reward="+quest.Reward)
        player.MsgChan <- models.Message{
            Category: "QUEST",
            Content:  "COMPLETED " + quest.ID,
        }
    default:
        speak.SendError(player.Conn, 400, "Usage: QUEST <ACCEPT|COMPLETE> <quest_id>")
    }
}