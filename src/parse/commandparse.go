package parse

import (
	"answer_protocol/src/game"
	"answer_protocol/src/logger"
	"answer_protocol/src/models"
	"answer_protocol/src/speakserver"
	"encoding/json"
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
	logger.Info("command received", "player", player.Name, "addr", logger.Addr(player.Conn), "cmd", command, "args", argument)
	if player.GetStatus() == "combat" {
		switch command {
		case "USE_ITEM", "DEFEND", "FLEE", "STATUS", "REQ":
			if argument != "" {
				speak.SendErr(player.Conn, speak.ErrUnexpectedArgument)
				return
			}
			if command == "USE_ITEM" {
				game.Attack(player, h)
				return
			}
			if command == "DEFEND" {
				game.Defend(player, h)
				return
			}
			if command == "FLEE" {
				game.Flee(player, h)
				return
			}
			if command == "STATUS" {
				game.ShowStatus(player)
				return
			}
			if command == "REQ" {
				return
			}
		default:
			speak.SendErr(player.Conn, speak.ErrCommandInCombat)
			return
		}
	}
	switch command {
	case "LOOK", "INVENTORY", "STATUS", "QUESTS", "WHO", "QUIT":
		if argument != "" {
			speak.SendErr(player.Conn, speak.ErrUnexpectedArgument)
			return
		}
		if command == "LOOK" {
			game.ShowRoom(player, h)
			return
		}
		if command == "INVENTORY" {
			game.ShowInventory(player)
			return
		}
		if command == "STATUS" {
			game.ShowStatus(player)
			return
		}
		if command == "QUESTS" {
			game.ShowQuest(player)
			return
		}
		if command == "WHO" {
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
			speak.SendErr(player.Conn, speak.ErrMissingArgument)
			return
		}
		argument = strings.ToUpper(argument)
		switch argument {
		case "NORTH", "SOUTH", "WEST", "EAST":
			game.MapMove(player, argument, h)
			return
		default:
			speak.SendErr(player.Conn, speak.ErrNoExit)
		}
	case "CHAT":
		if argument == "" {
			speak.SendErr(player.Conn, speak.ErrMissingArgument)
			return
		}
		partsChat := strings.SplitN(argument, " ", 2)
		if len(partsChat) < 2 {
			speak.SendErr(player.Conn, speak.ErrMalformedCommand)
			return
		}
		parseChat(partsChat, player, h)
	case "TAKE":
		if argument == "" {
			speak.SendErr(player.Conn, speak.ErrMissingArgument)
			return
		}
		game.TakeItem(player, argument, h)
	case "DROP":
		if argument == "" {
			speak.SendErr(player.Conn, speak.ErrMissingArgument)
			return
		}
		game.DropItem(player, argument, h)
	case "TALK":
		if argument == "" {
			speak.SendErr(player.Conn, speak.ErrMissingArgument)
			return
		}
		game.TalkNpc(player, argument)
	case "QUEST":
		if argument == "" {
			speak.SendErr(player.Conn, speak.ErrMissingArgument)
			return
		}
		partsQuest := strings.SplitN(argument, " ", 2)
		if len(partsQuest) < 2 {
			speak.SendErr(player.Conn, speak.ErrMalformedCommand)
			return
		}
		game.ManageQuest(player, h, partsQuest[0], partsQuest[1])
	case "ATTACK":
		if argument == "" {
			speak.SendErr(player.Conn, speak.ErrMissingArgument)
			return
		}
		game.StartAttack(player, argument, h)
	case "GROUP":
		if argument == "" {
			speak.SendErr(player.Conn, speak.ErrMissingArgument)
			return
		}
		partsGroup := strings.Split(argument, " ")
		parseGroup(partsGroup, player, h)
	case "REQ":
		player.Conn.Write([]byte{0x03})
		if player.Room != nil {
			player.Room.Mu.RLock()
			for i, npc := range player.Room.Npcs {
				if npc == nil {
					logger.Warn("nil npc pointer", "index", i, "room", player.Room.Id)
					continue
				}
			}
			player.Room.Mu.RUnlock()
		}
		state := models.WorldStateResponse{
			RoomItems:    player.GetCurrentRoomItemIDs(),
			RoNpcsTalk:   player.GetCurrentRoomNpcIDsTalk(),
			RoNpcsHostil: player.GetCurrentRoomNpcIDsHostil(),
			Inventory:    player.GetInventoryItemIDs(),
			PlayerQuests: player.GetPlayerQuestsList(),
			NpcQuests:    player.GetRoomNpcQuests(),
		}
		json.NewEncoder(player.Conn).Encode(state)
		return
	default:
		speak.SendErr(player.Conn, speak.ErrUnknownCommand)
		return
	}
}
