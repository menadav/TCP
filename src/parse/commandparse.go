package parse

import (
	"answer_protocol/src/models"
	"answer_protocol/src/speakserver"
	"answer_protocol/src/utils"
	"fmt"
	"strings"
	"net"
)

func ParseCommandCli(line string, player *models.Player, h *models.Hub) {
	line = strings.TrimSpace(line)
	if line == "" {
		return
	}
	parts := strings.SplitN(line, " ", 2)
	command := strings.ToUpper(parts[0])
	var argument string
	if len(parts) > 1 {
		argument = parts[1]
	}
	switch command {
	case "LOOK", "INVENTORY", "STATUS", "QUESTS", "WHO", "QUIT":
		if argument != "" {
			speak.SendError(player.Conn, 400, "Only command, no arguments allowed")
			return
		}
		if command == "QUIT" {
			h.Unregister <- player
			speak.SendSuccess(player.Conn, "bye")
			return
		}
	case "MOVE":
		if argument == "" {
			speak.SendError(player.Conn, 400, "Move requires a destination")
			return
		}
	case "CHAT":
		if argument == "" {
			speak.SendError(player.Conn, 400, "Chat requires a scope and a message")
			return
		}
		partsChat := strings.SplitN(argument, " ", 2)
		if len(partsChat) < 2 {
			speak.SendError(player.Conn, 400, "Chat format invalid. Use: CHAT <SCOPE> <MESSAGE>")
			return
		}
		parseChat(partsChat, player, h)
	case "TAKE", "DROP", "TALK", "ATTACK", "QUEST":
		if argument == "" {
			speak.SendError(player.Conn, 400, "This command requires an argument")
			return
		}
	case "GROUP":
		if argument == "" {
			return
		}
		partsGroup := strings.Split(argument, " ")
		parseGroup(partsGroup, player, h)
	default:
		speak.SendError(player.Conn, 400, "Unknown command")
		return
	}
}

func parseChat(partsChat []string, player *models.Player, h *models.Hub) {
	var msg models.Message

	scopeStr := strings.ToUpper(partsChat[0])
	text := partsChat[1]
	switch scopeStr {
	case "GLOBAL":
		msg = models.Message{
			Scope:   models.ScopeGlobal,
			Filter:  "",
			Category: "GLOBAL",
			Content: fmt.Sprintf("CHAT %s %s", player.Name, text),
		}
	case "ROOM":
		msg = models.Message{
			Scope:   models.ScopeRoom,
			Filter:  player.Room,
			Category: "ROOM",
			Content: fmt.Sprintf("CHAT %s %s", player.Name, text),
		}
	case "GROUP":
		if player.Group == "" {
			speak.SendError(player.Conn, 403, "You are not in a group")
			return
		}
		msg = models.Message{
			Scope:    models.ScopeGroup,
			Filter:   player.Group,
			Category: "GROUP",
			Content:  fmt.Sprintf("CHAT %s %s", player.Name, text),
		}
	default:
		speak.SendError(player.Conn, 400, "Unknown chat scope. Use GLOBAL, ROOM, or GROUP")
		return
	}
	h.Broadcast <- msg
}

func parseGroup(partsGroup []string, player *models.Player, h *models.Hub){
	action := strings.ToUpper(partsGroup[0])
	switch action{
	case "CREATE":
		if player.Group != "" {
			speak.SendError(player.Conn, 403, "Already in a group")
			return
		}
		groupID := "grp_" + player.Name
		newGroup := &models.Group{
			Id:			groupID,
			Leader:		player,
			Members:	make(map[net.Conn]*models.Player),
		}
		newGroup.AddMember(player.Conn, player)
		player.Group = groupID
		h.Groups[groupID] = newGroup
		speak.SendSuccess(player.Conn, fmt.Sprintf("group=%s", groupID))
	case "INVITE":
		if player.Group == "" {
			speak.SendError(player.Conn, 403, "You are not in a group")
			return
		}
		actualGroup, exist := h.Groups[player.Group]
		if !exist {
            speak.SendError(player.Conn, 404, "Group not found in server records")
            return
        }
		if actualGroup.Leader != player {
            speak.SendError(player.Conn, 403, "Only the leader can invite players")
            return
        }
		if len(partsGroup) < 2 {
			speak.SendError(player.Conn, 403, "Missing username to invite")
			return
		}
		targetUsername := partsGroup[1]
		if !utils.ExistName(h.Clients ,targetUsername) {
			speak.SendError(player.Conn, 404, "User not found")
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
			speak.SendError(player.Conn, 404, "User log out or not found")
			return
		}
		if targetPlayer.Group != "" {
			speak.SendError(player.Conn, 403, "User is already in another group")
			return
		}
		speak.SendEvent(targetPlayer.Conn, "GROUP INVITE", player.Name )
		speak.SendSuccess(player.Conn, "")
	case "JOIN":
		if player.Group != "" {
			speak.SendError(player.Conn, 403, "Already in a group. Leave current group first.")
			return
		}
		if len(partsGroup) < 2 {
			speak.SendError(player.Conn, 400, "Missing leader name to join")
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
			speak.SendError(player.Conn, 404, "No active group found with that leader")
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
            speak.SendError(player.Conn, 401, "NOT_IN_GROUP")
            return
        }
        groupID := player.Group
        if group, exist := h.Groups[groupID]; exist {
            All_group := group.RemoveMember(player.Conn)
            player.Group = ""
            speak.SendSuccess(player.Conn, "")
            if All_group == 0 || group.Leader == player {
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
		speak.SendError(player.Conn, 400, "Unknown group action. Use CREATE, INVITE, or JOIN")
		return
	}
}