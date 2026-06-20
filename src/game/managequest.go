package game

import (
	"answer_protocol/src/models"
	"answer_protocol/src/speakserver"
	"strings"
)

func ManageQuest(player *models.Player, h *models.Hub, action string, questID string) {
	action = strings.ToUpper(action)

	quest, exists := h.World.Quest[questID]
	if !exists {
		speak.SendErr(player.Conn, speak.ErrQuestNotFound)
		return
	}
	switch action {
	case "ACCEPT":
		var startItem *models.Item
		if quest.StartItem != "" {
			startItem = h.World.Items[quest.StartItem]
		}
		if e := player.AcceptQuest(quest, startItem); e != nil {
			speak.SendErr(player.Conn, *e)
			return
		}
		speak.SendSuccess(player.Conn, "quest_accepted="+quest.ID)
	case "COMPLETE":
		rewardItem := h.World.Items[quest.Reward]
		if e := player.CompleteQuest(quest, rewardItem); e != nil {
			speak.SendErr(player.Conn, *e)
			return
		}
		speak.SendSuccess(player.Conn, "quest_completed="+quest.ID+" reward="+quest.Reward)
		player.MsgChan <- models.Message{
			Category: "QUEST",
			Content:  "COMPLETED " + quest.ID,
		}
	default:
		speak.SendErr(player.Conn, speak.ErrInvalidArgument)
	}
}
