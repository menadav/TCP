package parse

import (
	"answer_protocol/src/models"
	"answer_protocol/src/speakserver"
	"answer_protocol/src/utils"
	"fmt"
	"net"
	"strings"
)

func parseGroup(partsGroup []string, player *models.Player, h *models.Hub) {
	action := strings.ToUpper(partsGroup[0])
	switch action {
	case "CREATE":
		if player.Group != "" {
			speak.SendErr(player.Conn, speak.ErrAlreadyInGroup)
			return
		}
		groupID := "grp_" + player.Name
		newGroup := &models.Group{
			Id:      groupID,
			Leader:  player,
			Members: make(map[net.Conn]*models.Player),
		}
		newGroup.AddMember(player.Conn, player)
		player.Group = groupID
		h.Groups[groupID] = newGroup
		speak.SendSuccess(player.Conn, fmt.Sprintf("group=%s", groupID))
	case "INVITE":
		if player.Group == "" {
			speak.SendErr(player.Conn, speak.ErrNotInGroup)
			return
		}
		actualGroup, exist := h.Groups[player.Group]
		if !exist {
			speak.SendErr(player.Conn, speak.ErrGroupNotFound)
			return
		}
		if actualGroup.Leader != player {
			speak.SendErr(player.Conn, speak.ErrNotGroupLeader)
			return
		}
		if len(partsGroup) < 2 {
			speak.SendErr(player.Conn, speak.ErrMissingArgument)
			return
		}
		targetUsername := partsGroup[1]
		if !utils.ExistName(h.Clients, targetUsername) {
			speak.SendErr(player.Conn, speak.ErrUserNotFound)
			return
		}
		var targetPlayer *models.Player
		for _, client := range h.Clients {
			if client.Name == targetUsername {
				targetPlayer = client
				break
			}
		}
		if targetPlayer == nil {
			speak.SendErr(player.Conn, speak.ErrUserNotFound)
			return
		}
		if targetPlayer.Group != "" {
			speak.SendErr(player.Conn, speak.ErrAlreadyInGroup)
			return
		}
		speak.SendEvent(targetPlayer.Conn, "GROUP INVITE", player.Name)
		speak.SendSuccess(player.Conn, "")
	case "JOIN":
		if player.Group != "" {
			speak.SendErr(player.Conn, speak.ErrAlreadyInGroup)
			return
		}
		if len(partsGroup) < 2 {
			speak.SendErr(player.Conn, speak.ErrMissingArgument)
			return
		}
		leaderName := partsGroup[1]
		var targetGroup *models.Group

		for _, g := range h.Groups {
			if g.Leader.Name == leaderName {
				targetGroup = g
				break
			}
		}
		if targetGroup == nil {
			speak.SendErr(player.Conn, speak.ErrGroupNotFound)
			return
		}
		h.Broadcast <- models.Message{
			Scope:    models.ScopeGroup,
			Filter:   targetGroup.Id,
			Category: "GROUP",
			Content:  "JOIN " + player.Name,
		}
		targetGroup.AddMember(player.Conn, player)
		player.Group = targetGroup.Id
		speak.SendSuccess(player.Conn, fmt.Sprintf("group=%s", targetGroup.Id))
	case "LEAVE":
		if player.Group == "" {
			speak.SendErr(player.Conn, speak.ErrNotInGroup)
			return
		}
		groupID := player.Group
		if group, exist := h.Groups[groupID]; exist {
			isLeader := group.Leader == player
			All_group := group.RemoveMember(player.Conn)
			player.Group = ""
			speak.SendSuccess(player.Conn, "")
			if isLeader {
				for _, member := range group.Members {
					speak.SendEvent(member.Conn, "GROUP", "DISBANDED")
					member.Group = ""
				}
				h.Broadcast <- models.Message{
					Scope:    models.ScopeGroup,
					Filter:   groupID,
					Category: "GROUP",
					Content:  "DISBAND " + player.Name,
				}
				delete(h.Groups, groupID)
			} else if All_group == 0 {
				h.Broadcast <- models.Message{
					Scope:    models.ScopeGroup,
					Filter:   groupID,
					Category: "GROUP",
					Content:  "LEAVE " + player.Name,
				}
				delete(h.Groups, groupID)
			} else {
				h.Broadcast <- models.Message{
					Scope:    models.ScopeGroup,
					Filter:   groupID,
					Category: "GROUP",
					Content:  "LEAVE " + player.Name,
				}
			}
		}
	default:
		speak.SendErr(player.Conn, speak.ErrInvalidArgument)
		return
	}
}
