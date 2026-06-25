package gui

import (
	"answer_protocol/src/models"
	"bufio"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

type compactTheme struct {
	fyne.Theme
}

func (c compactTheme) Size(name fyne.ThemeSizeName) float32 {
	switch name {
	case theme.SizeNameText:
		return 11
	case theme.SizeNamePadding:
		return 2
	}
	return c.Theme.Size(name)
}

type GameUI struct {
	App           fyne.App
	Window        fyne.Window
	MudConsole    *widget.RichText
	ChatView      *widget.RichText
	RoomPanel     *widget.RichText
	CountersLabel *widget.Label
	ChatInput     *widget.Entry
	ActionsPanel  *fyne.Container
	CurrentPlayer *models.Player
	History       string
	ChatHistory   string
	whoLogPending bool
	currentMenu   string
	LastState     models.WorldStateResponse
}

func Start(app fyne.App, player *models.Player) {
	app.Settings().SetTheme(compactTheme{theme.DefaultTheme()})
	if player.Room == nil {
		player.Room = &models.Room{Exist: make(map[string]string)}
	}

	ui := &GameUI{
		App:           app,
		Window:        app.NewWindow("TAP MUD"),
		MudConsole:    widget.NewRichTextFromMarkdown("## Log\n"),
		ChatView:      widget.NewRichTextFromMarkdown("## Chat\n"),
		RoomPanel:     widget.NewRichTextFromMarkdown("### Room\n"),
		CountersLabel: widget.NewLabel("Players  room: -  |  server: -"),
		ChatInput:     widget.NewEntry(),
		CurrentPlayer: player,
		History:       "## Log\n",
		ChatHistory:   "## Chat\n",
		currentMenu:   "exploration",
	}
	ui.ChatInput.SetPlaceHolder("Type command or chat...")
	ui.MudConsole.Wrapping = fyne.TextWrapWord
	ui.ChatView.Wrapping = fyne.TextWrapWord
	ui.RoomPanel.Wrapping = fyne.TextWrapWord
	go ui.listenServer()
	ui.setupLayout()
	ui.Window.Resize(fyne.NewSize(820, 470))
	ui.Window.ShowAndRun()
}

func (ui *GameUI) listenServer() {
	reader := bufio.NewReader(ui.CurrentPlayer.Conn)
	for {
		opcode, err := reader.ReadByte()
		if err != nil {
			return
		}

		if opcode == 0x03 {
			var state models.WorldStateResponse
			if err := json.NewDecoder(reader).Decode(&state); err == nil {
				ui.LastState = state
				newInventory := make([]*models.Item, len(state.Inventory))
				for i, iv := range state.Inventory {
					newInventory[i] = &models.Item{ID: iv.ID, Name: iv.Name}
				}
				ui.CurrentPlayer.Inventory = newInventory
				ui.CurrentPlayer.Room.Mu.Lock()
				ui.CurrentPlayer.Room.Items = make([]*models.Item, len(state.RoomItems))
				for i, iv := range state.RoomItems {
					ui.CurrentPlayer.Room.Items[i] = &models.Item{ID: iv.ID, Name: iv.Name}
				}
				ui.CurrentPlayer.Room.Exist = make(map[string]string)
				for _, id := range state.RoNpcsTalk {
					ui.CurrentPlayer.Room.Exist[id] = "TALK"
				}
				for _, id := range state.RoNpcsHostil {
					ui.CurrentPlayer.Room.Exist[id] = "HOSTILE"
				}
				ui.CurrentPlayer.Room.Mu.Unlock()
				ui.updateRoomPanel(state)
				ui.refreshCurrentMenu()

			}
			continue
		}
		msg, _ := reader.ReadString('\n')
		if strings.Contains(msg, "ITEMS_CHANGED") {
			fmt.Fprintf(ui.CurrentPlayer.Conn, "REQ\n")
			continue
		}
		line := strings.TrimSpace(string(opcode) + msg)
		if strings.HasPrefix(line, "OK who=") {
			var who models.WhoResponse
			if err := json.Unmarshal([]byte(strings.TrimPrefix(line, "OK who=")), &who); err == nil {
				ui.CountersLabel.SetText(fmt.Sprintf("Players  room: %d  |  server: %d", len(who.Room), who.Server))
				if ui.whoLogPending {
					ui.whoLogPending = false
					ui.History += fmt.Sprintf("\n\n**WHO** - room (%d): %s - server total: %d", len(who.Room), strings.Join(who.Room, ", "), who.Server)
					ui.MudConsole.ParseMarkdown(ui.History)
					ui.MudConsole.Refresh()
				}
			}
			continue
		}
		if isChatLine(line) {
			ui.ChatHistory += "\n\n" + line
			ui.ChatView.ParseMarkdown(ui.ChatHistory)
			ui.ChatView.Refresh()
		} else {
			ui.History += "\n\n" + line
			ui.MudConsole.ParseMarkdown(ui.History)
			ui.MudConsole.Refresh()
			if strings.Contains(line, "PRESENCE") || strings.HasPrefix(line, "OK connected") {
				fmt.Fprintf(ui.CurrentPlayer.Conn, "WHO\n")
			}
		}
	}
}

func isChatLine(line string) bool {
	fields := strings.Fields(line)
	return len(fields) >= 3 && fields[0] == "EVT" && fields[2] == "CHAT"
}

func (ui *GameUI) updateRoomPanel(state models.WorldStateResponse) {
	order := []string{"NORTH", "SOUTH", "EAST", "WEST"}
	var parts []string
	for _, dir := range order {
		if target, ok := state.RoomExits[dir]; ok {
			parts = append(parts, dir+" -> "+target)
		}
	}
	exits := strings.Join(parts, ", ")
	if exits == "" {
		exits = "none"
	}
	md := fmt.Sprintf("### %s\n%s\n\n**Exits:** %s", state.RoomName, state.RoomDesc, exits)
	ui.RoomPanel.ParseMarkdown(md)
	ui.RoomPanel.Refresh()
}

func (ui *GameUI) sendCommand(cmd string) {
	cmd = strings.TrimSpace(cmd)

	if cmd == "" {
		return
	}
	if strings.ToUpper(cmd) == "WHO" {
		ui.whoLogPending = true
	}
	fmt.Fprintf(ui.CurrentPlayer.Conn, "%s\n", cmd)
	if !strings.HasPrefix(cmd, "CHAT") && !strings.HasPrefix(cmd, "LOGIN") && cmd != "REQ" {
		time.Sleep(100 * time.Millisecond)
		fmt.Fprintf(ui.CurrentPlayer.Conn, "REQ\n")
	}
}

func (ui *GameUI) setupLayout() {
	ui.ActionsPanel = container.NewVBox()
	ui.showExplorationMenu()

	chatPrefixes := container.NewHBox(
		widget.NewButton("Global", func() { ui.ChatInput.SetText("CHAT GLOBAL ") }),
		widget.NewButton("Room", func() { ui.ChatInput.SetText("CHAT ROOM ") }),
		widget.NewButton("Group", func() { ui.ChatInput.SetText("CHAT GROUP ") }),
	)

	movePanel := container.NewGridWithColumns(3,
		layout.NewSpacer(),
		widget.NewButton("N", func() { ui.sendCommand("MOVE NORTH") }),
		layout.NewSpacer(),
		widget.NewButton("W", func() { ui.sendCommand("MOVE WEST") }),
		layout.NewSpacer(),
		widget.NewButton("E", func() { ui.sendCommand("MOVE EAST") }),
		layout.NewSpacer(),
		widget.NewButton("S", func() { ui.sendCommand("MOVE SOUTH") }),
		layout.NewSpacer(),
	)

	btnSend := widget.NewButton("Send", func() {
		if ui.ChatInput.Text != "" {
			ui.sendCommand(ui.ChatInput.Text)
			ui.ChatInput.SetText("")
		}
	})

	bottom := container.NewVBox(
		chatPrefixes,
		widget.NewSeparator(),
		movePanel,
		widget.NewSeparator(),
		container.NewBorder(nil, nil, nil, btnSend, ui.ChatInput),
	)

	chatScroll := container.NewScroll(ui.ChatView)
	logScroll := container.NewScroll(ui.MudConsole)
	chatScroll.SetMinSize(fyne.NewSize(200, 150))
	logScroll.SetMinSize(fyne.NewSize(200, 150))
	chatPanel := container.NewBorder(widget.NewLabelWithStyle("Chat", fyne.TextAlignCenter, fyne.TextStyle{Bold: true}), nil, nil, nil, chatScroll)
	logPanel := container.NewBorder(widget.NewLabelWithStyle("Log", fyne.TextAlignCenter, fyne.TextStyle{Bold: true}), nil, nil, nil, logScroll)
	roomScroll := container.NewScroll(ui.RoomPanel)
	roomScroll.SetMinSize(fyne.NewSize(180, 150))
	roomPanel := container.NewBorder(widget.NewLabelWithStyle("Room", fyne.TextAlignCenter, fyne.TextStyle{Bold: true}), nil, nil, nil, roomScroll)
	consoleSplit := container.NewHSplit(chatPanel, logPanel)
	consoleSplit.SetOffset(0.5)
	leftSide := container.NewHSplit(roomPanel, consoleSplit)
	leftSide.SetOffset(0.34)

	main := container.NewHSplit(leftSide, ui.ActionsPanel)
	main.SetOffset(0.72)
	topBar := container.NewHBox(ui.CountersLabel)
	ui.Window.SetContent(container.NewBorder(topBar, bottom, nil, nil, main))
}

func (ui *GameUI) showCombatMenu(targetNpc string) {
	ui.currentMenu = "combat"
	ui.ActionsPanel.Objects = []fyne.CanvasObject{
		widget.NewLabel("=== COMBAT: " + targetNpc + " ==="),
		widget.NewButton("USE_ITEM", func() {
			ui.sendCommand("USE_ITEM")
		}),
		widget.NewButton("DEFEND", func() {
			ui.sendCommand("DEFEND")
		}),
		widget.NewButton("FLEE", func() {
			ui.sendCommand("FLEE")
		}),
		widget.NewButton("STATUS", func() {
			ui.sendCommand("STATUS")
		}),
	}
	ui.ActionsPanel.Refresh()
}

func (ui *GameUI) groupButtons() fyne.CanvasObject {
	return container.NewGridWithColumns(2,
		widget.NewButton("Create", func() {
			ui.sendCommand("GROUP CREATE")
		}),
		widget.NewButton("Leave", func() {
			ui.sendCommand("GROUP LEAVE")
		}),
		widget.NewButton("Join...", func() {
			ui.ChatInput.SetText("GROUP JOIN ")
			ui.ChatInput.FocusGained()
		}),
		widget.NewButton("Invite...", func() {
			ui.ChatInput.SetText("GROUP INVITE ")
			ui.ChatInput.FocusGained()
		}),
	)
}

func (ui *GameUI) showExplorationMenu() {
	ui.currentMenu = "exploration"

	groupButtons := container.NewGridWithColumns(2,
		widget.NewButton("Create", func() { ui.sendCommand("GROUP CREATE") }),
		widget.NewButton("Leave", func() { ui.sendCommand("GROUP LEAVE") }),
		widget.NewButton("Join...", func() { ui.ChatInput.SetText("GROUP JOIN ") }),
		widget.NewButton("Invite...", func() { ui.ChatInput.SetText("GROUP INVITE ") }),
	)

	ui.ActionsPanel.Objects = []fyne.CanvasObject{
		widget.NewLabel("=== ACTIONS ==="),
		widget.NewButton("LOOK", func() { ui.sendCommand("LOOK") }),
		widget.NewButton("INVENTORY", func() {
			ui.sendCommand("INVENTORY")
			ui.showInventoryMenu()
		}),
		widget.NewButton("STATUS", func() { ui.sendCommand("STATUS") }),
		widget.NewButton("WHO", func() { ui.sendCommand("WHO") }),
		widget.NewButton("QUESTS", func() { ui.sendCommand("QUESTS") }),

		widget.NewLabel("=== GROUP ==="),
		groupButtons,
		widget.NewLabel("=== INTERACT ==="),
		widget.NewButton("ATTACK...", func() { ui.showTargetSelection("ATTACK") }),
		widget.NewButton("TALK...", func() { ui.showTargetSelection("TALK") }),
		widget.NewButton("TAKE...", func() { ui.showTargetSelection("TAKE") }),
		widget.NewButton("DROP...", func() { ui.showTargetSelection("DROP") }),
		widget.NewButton("QUEST ACCEPT...", func() { ui.showTargetSelection("QUEST_ACCEPT") }),
		widget.NewButton("QUEST COMPLETE...", func() { ui.showTargetSelection("QUEST_COMPLETE") }),
		widget.NewLabel("---"),
		widget.NewButton("QUIT", func() {
			ui.sendCommand("QUIT")
			ui.Window.Close()
		}),
	}
	ui.ActionsPanel.Refresh()
}

func (ui *GameUI) removeLocalItem(itemToRemove *models.Item) {
	newInv := []*models.Item{}
	for _, item := range ui.CurrentPlayer.Inventory {
		if item != itemToRemove {
			newInv = append(newInv, item)
		}
	}
	ui.CurrentPlayer.Inventory = newInv
}

func itemLabel(item *models.Item) string {
	if item == nil {
		return ""
	}
	if item.Name != "" {
		return item.Name
	}
	return item.ID
}

func (ui *GameUI) showTargetSelection(action string) {
	ui.currentMenu = action
	var objs []fyne.CanvasObject
	objs = append(objs, widget.NewLabel(fmt.Sprintf("=== %s ===", action)))

	if action == "QUEST_ACCEPT" {
		for _, quest := range ui.LastState.NpcQuests {
			qID := quest.ID
			qTitle := quest.Title
			objs = append(objs, widget.NewButton("Accept: "+qTitle, func() {
				ui.sendCommand("QUEST ACCEPT " + qID)
			}))
		}
	} else if action == "QUEST_COMPLETE" {
		for _, q := range ui.LastState.PlayerQuests {
			qID := q.QuestID
			progress := q.Progress

			objs = append(objs, widget.NewButton("Complete: "+qID+" ("+progress+")", func() {
				ui.sendCommand("QUEST COMPLETE " + qID)
			}))
		}
	} else {
		ui.CurrentPlayer.Room.Mu.RLock()

		if action == "ATTACK" || action == "TALK" {
			for id, kind := range ui.CurrentPlayer.Room.Exist {
				if (action == "ATTACK" && kind == "HOSTILE") || (action == "TALK" && kind == "TALK") {
					id := id
					objs = append(objs, widget.NewButton(action+" "+id, func() {
						ui.sendCommand(action + " " + id)
						if action == "ATTACK" {
							ui.showCombatMenu(id)
						}
					}))
				}
			}
		} else if action == "TAKE" {
			if len(ui.CurrentPlayer.Room.Items) == 0 {
				objs = append(objs, widget.NewLabel("No items here."))
			}
			for _, item := range ui.CurrentPlayer.Room.Items {
				id := item.ID
				objs = append(objs, widget.NewButton("Take "+itemLabel(item), func() { ui.sendCommand("TAKE " + id) }))
			}
		} else if action == "DROP" {
			items := ui.CurrentPlayer.Inventory
			if len(items) == 0 {
				objs = append(objs, widget.NewLabel("Inventory is empty."))
			} else {
				for _, item := range items {
					if item == nil {
						continue
					}
					targetItem := item
					objs = append(objs, widget.NewButton("Drop "+itemLabel(targetItem), func() {
						ui.sendCommand("DROP " + targetItem.ID)
						ui.removeLocalItem(targetItem)
						ui.showTargetSelection("DROP")
					}))
				}
			}
		}
		ui.CurrentPlayer.Room.Mu.RUnlock()
	}

	objs = append(objs, widget.NewButton("Back", ui.showExplorationMenu))
	ui.ActionsPanel.Objects = objs
	ui.ActionsPanel.Refresh()
}

func (ui *GameUI) showInventoryMenu() {
	ui.currentMenu = "inventory"
	var objs []fyne.CanvasObject
	objs = append(objs, widget.NewLabel("=== INVENTORY ==="))

	items := ui.CurrentPlayer.Inventory
	if len(items) == 0 {
		objs = append(objs, widget.NewLabel("Your inventory is empty."))
	} else {
		for _, item := range items {
			if item == nil {
				continue
			}
			currentID := item.ID
			buttonText := "Drop " + itemLabel(item)
			objs = append(objs, widget.NewButton(buttonText, func() {
				ui.sendCommand("DROP " + currentID)
			}))
		}
	}

	objs = append(objs, widget.NewButton("Back", ui.showExplorationMenu))
	ui.ActionsPanel.Objects = objs
	ui.ActionsPanel.Refresh()
}

func (ui *GameUI) refreshCurrentMenu() {
	ui.CurrentPlayer.Mu.RLock()
	isFighting := ui.CurrentPlayer.CombatNpc != ""
	target := ui.CurrentPlayer.CombatNpc
	ui.CurrentPlayer.Mu.RUnlock()

	if isFighting {
		ui.showCombatMenu(target)
		return
	}
	switch ui.currentMenu {
	case "exploration":
		ui.showExplorationMenu()
	case "ATTACK", "TALK", "TAKE", "DROP":
		ui.showTargetSelection(ui.currentMenu)
	case "QUEST_ACCEPT", "QUEST_COMPLETE":
		ui.showTargetSelection(ui.currentMenu)
	case "inventory":
		ui.showInventoryMenu()
	default:
		ui.showExplorationMenu()
	}
}
