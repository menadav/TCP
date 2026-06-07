package parse

import (
	"answer_protocol/src/models"
	"answer_protocol/src/speakserver"
	"answer_protocol/src/game"
	"strings"
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
		if command == "LOOK"{
			game.ShowRoom(player, h)
			return
		}
		if command == "INVENTORY"{
			game.ShowInventory(player)
			return
    	}
		if command == "STATUS"{
			game.ShowStatus(player)
			return
		}
		if command == "QUESTS"{
			game.ShowQuest(player)
			return
		}
		if command == "WHO"{
			game.ShowWho(player, h)
			return
		}
		if command == "QUIT" {
			speak.SendSuccess(player.Conn, "bye")
			player.Conn.Close()
			return
		}
	case "MOVE":
		if argument == "" {
			speak.SendError(player.Conn, 400, "Move requires a destination")
			return
		}
		argument = strings.ToUpper(argument)
		switch argument {
		case "NORTH", "SOUTH", "WEST", "EAST":
			game.MapMove(player, argument, h)
			return
		default:
			speak.SendError(player.Conn, 301, "NO_EXIT")
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
	case "TAKE":
		if argument == "" {
			speak.SendError(player.Conn, 400, "TAKE requires an item name")
			return
		}
		game.TakeItem(player, argument)
	case "DROP":
		if argument == "" {
			speak.SendError(player.Conn, 400, "DROP requires an item name")
			return
		}
		game.DropItem(player, argument)
	case "TALK":
		if argument == "" {
			speak.SendError(player.Conn, 400, "TALK requires an NPC name")
			return
		}
		game.TalkNpc(player, argument)
	case "ATTACK":
		if argument == "" {
			speak.SendError(player.Conn, 400, "ATTACK requires an NPC name")
			return
		}
		game.AttackNpc(player, argument, h)
	case "QUEST":
		if argument == "" {
			speak.SendError(player.Conn, 400, "QUEST requires a quest ID")
			return
		}
		game.HandleQuest(player, argument, h.World)
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