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
    "fyne.io/fyne/v2/widget"
)
type GameUI struct {
	App           fyne.App
	Window        fyne.Window
	MudConsole    *widget.RichText
	ChatInput     *widget.Entry
	ActionsPanel  *fyne.Container
	CurrentPlayer *models.Player
	History       string
	currentMenu   string
	LastState     models.WorldStateResponse
}

func Start(app fyne.App, player *models.Player) {
    if player.Room == nil {
        player.Room = &models.Room{Exist: make(map[string]string)}
    }

    ui := &GameUI{
        App:           app,
        Window:        app.NewWindow("TAP MUD"),
        MudConsole:    widget.NewRichTextFromMarkdown("## Welcome to TAP\n"),
        ChatInput:     widget.NewEntry(),
        CurrentPlayer: player,
        History:       "## Welcome to TAP\n",
        currentMenu:   "exploration",
    }
    ui.ChatInput.SetPlaceHolder("Type command or chat...")
    ui.MudConsole.Wrapping = fyne.TextWrapWord
    go ui.listenServer()
    ui.setupLayout()
    ui.Window.Resize(fyne.NewSize(900, 600))
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
				for i, id := range state.Inventory {
					newInventory[i] = &models.Item{ID: id}
				}
				ui.CurrentPlayer.Inventory = newInventory
				ui.CurrentPlayer.Room.Mu.Lock()
				ui.CurrentPlayer.Room.Items = make([]*models.Item, len(state.RoomItems))
				for i, id := range state.RoomItems {
					ui.CurrentPlayer.Room.Items[i] = &models.Item{ID: id}
				}
				ui.CurrentPlayer.Room.Exist = make(map[string]string)
				for _, id := range state.RoNpcsTalk {
					ui.CurrentPlayer.Room.Exist[id] = "TALK"
				}
				for _, id := range state.RoNpcsHostil {
					ui.CurrentPlayer.Room.Exist[id] = "HOSTILE"
				}
				ui.CurrentPlayer.Room.Mu.Unlock()
				ui.refreshCurrentMenu()

			}
			continue
		}
		msg, _ := reader.ReadString('\n')
		if strings.Contains(msg, "ITEMS_CHANGED") {
			fmt.Fprintf(ui.CurrentPlayer.Conn, "REQ\n")
			continue
		}
		ui.History += "\n\n" + strings.TrimSpace(string(opcode)+msg)
		ui.MudConsole.ParseMarkdown(ui.History)
		ui.MudConsole.Refresh()
	}
}

func (ui *GameUI) sendCommand(cmd string) {
	cmd = strings.TrimSpace(cmd)

	if cmd == "" {
		return
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
    
    consoleScroll := container.NewScroll(ui.MudConsole)
    consoleScroll.SetMinSize(fyne.NewSize(400, 400))

    main := container.NewHSplit(consoleScroll, ui.ActionsPanel)
    main.SetOffset(0.7)    
    ui.Window.SetContent(container.NewBorder(nil, bottom, nil, nil, main))
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
						ui.showCombatMenu(id)
					}))
				}
			}
		} else if action == "TAKE" {
			if len(ui.CurrentPlayer.Room.Items) == 0 {
				objs = append(objs, widget.NewLabel("No items here."))
			}
			for _, item := range ui.CurrentPlayer.Room.Items {
				id := item.ID
				objs = append(objs, widget.NewButton("Take "+id, func() { ui.sendCommand("TAKE " + id) }))
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
					objs = append(objs, widget.NewButton("Drop "+targetItem.ID, func() {
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
			buttonText := "Drop " + currentID
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
