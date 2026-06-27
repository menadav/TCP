# 🎮 TAP — The Answer Protocol

> A real-time multiplayer text adventure (MUD) engine built on a custom TCP protocol, written in Go.

![TAP World Workflow](assets/mud.png)

[![Go](https://img.shields.io/badge/Go-1.18%2B-00ADD8?style=flat&logo=go)](https://go.dev/)
[![Fyne](https://img.shields.io/badge/GUI-Fyne-informational?style=flat)](https://fyne.io/)
[![Build](https://img.shields.io/badge/build-make-brightgreen?style=flat)]()
[![License](https://img.shields.io/badge/license-MIT-blue.svg)](LICENSE)

---

## 📖 What is this?

**TAP — The Answer Protocol** is a shared-world, retro text adventure where multiple players connect simultaneously to a single server and explore rooms, chat, form groups, fight NPCs and complete quests — all in real time over a custom line-based TCP protocol.

The project ships three components: a **TCP server**, a **CLI client**, and a **GUI client** built with [Fyne](https://fyne.io/). All world state lives in memory; no external database required.

---

## ✨ Features

- **Real-time multiplayer** — unlimited concurrent TCP connections; all state changes are broadcast live to every client
- **Dual clients** — lightweight terminal CLI and a Fyne GUI with live player counters, D-pad movement and separated chat/log views
- **Custom protocol** — versioned, line-based `OK` / `ERR` / `EVT` framing with a `CONNECT` handshake
- **Turn-based combat** — deterministic damage, `DEFEND`/`FLEE` mechanics, NPC counter-attacks and auto-respawn
- **Quest system** — `collect`, `kill`, `explore` and `deliver` quest types tracked per player
- **Group system** — create, invite, join and leave named groups with group-scoped chat and events
- **Structured JSON logging** — RFC3339 timestamps, levelled entries (`INFO`/`WARN`/`ERROR`) and abuse detection
- **Data-driven world** — rooms, items, NPCs and quests defined in `data.yaml`; no hardcoded world state

---

## 🏗️ Architecture

The server uses a **central hub + one goroutine per client** model. All global state mutations are serialised through a single `Hub` goroutine; per-entity `sync.RWMutex` locks protect `Room`, `Player` and `Group` structs.

```
cmd/server ──► network.ClientAtender        (1 goroutine / client)
                     │
                     ├─ network.Authentication   (CONNECT handshake)
                     ├─ parse.ParseCommandCli ──► src/game/*  (command handlers)
                     └─ player.ListenMsg          (1 goroutine / client → EVT events)
                     ▲
   models.Hub.Run ───┘   (single goroutine: Register / Unregister / Broadcast)
```

Key design decisions:

- **No shared mutable global state** — the client map lives exclusively inside `Hub.Run`; connections communicate via typed channels
- **Buffered per-player `MsgChan`** — asynchronous event delivery keeps clients responsive without blocking command processing
- **Inline command dispatch** — `parse.ParseCommandCli` is a `switch` over a small, fixed command set; readable and easy to audit

---

## 🚀 Quick Start

### 1. Install dependencies

```bash
make install
```

Fetches all Go module dependencies via `go mod download`.

### 2. Build all binaries

```bash
make build
```

Compiles server, CLI client and GUI client into `./bin`.

### 3. Run a multiplayer session

```bash
make run-server       # terminal 1 — starts the TCP server on :8080
make run-client       # terminal 2 — connect a CLI client
make run-client-gui   # terminal 3 — connect a GUI client
```

---

## ⚙️ Make Targets

| Target | Description |
|--------|-------------|
| `make install` | Download module dependencies |
| `make build` | Compile server, CLI and GUI into `./bin` |
| `make run-server` | Start the TCP server on `:8080` |
| `make run-client` | Start the CLI client |
| `make run-client-gui` | Start the GUI client |
| `make lint` | `gofmt` check + `go vet` |
| `make fmt-fix` | Reformat source with `gofmt -w` |
| `make clean` | Remove built binaries |

---

## 🌐 Protocol

Transport: **TCP, UTF-8, one message per line** (`\n` terminated).

### Handshake

```
← OK hello proto=1
→ CONNECT <name>          # 3–12 letters, unique across the session
← OK welcome name=<name>
```

### Message prefixes

| Prefix | Meaning |
|--------|---------|
| `OK <payload>` | Success — payload is `key=value` or a JSON object |
| `ERR <code> <SYMBOL>` | Failure — see [Error Reference](#-error-reference) |
| `EVT <category> <data>` | Async server push (chat, presence, combat…) |

### Non-standard extensions

| Command | Behaviour |
|---------|-----------|
| `REQ` | Returns a JSON snapshot (`WorldStateResponse`) of the current room — used by the GUI to refresh its view |
| `DEFEND` / `FLEE` | Combat-only commands |
| `WHO` | Replies `OK who={"room":[names],"server":N}` — live player counters |
| `INVENTORY` | Replies with a JSON array of the player's items |

---

## 🗺️ World Design

The world is themed around the **Kanto region** and ships with **8 rooms**, **4 items**, **3 NPCs** and **2 quests**.

```
              +------------------+
              |  lavender_town   |
              | (Lavender Town)  |
              +--------+---------+
                   N   |   S
+---------------+  W   +------------------+
|     tower     |<-----|    crossroad     |
| (Pkmn Tower)  |----->|(Lavender Crossrd)|
+---------------+  E   +--------+---------+
                         S  |  N
              +--------------------------+
              |     silence_bridge       |  <- Snorlax blocks NORTH
              |    (Silence Bridge)      |
              +--------+-----------------+
                   N   |   S
              +--------+---------+   E   +----------------+
              |     route_11     |------>|   graveyard    |
              |   (Route 11)     |<------| (Graveyard)    |
              +--------+---------+   W   +----------------+
                   S   |   N
              +------------------+
              |      start       |
              | (Vermilion City) |   <- Spawn point
              +------------------+
```

### Rooms

| Room | Name | Notable |
|------|------|---------|
| `start` | Vermilion City | Spawn point; Club President (quest-giver) |
| `graveyard` | Graveyard | Thick Bone weapon |
| `route_11` | Route 11 | — |
| `guard_house` | Guard House | — |
| `silence_bridge` | Silence Bridge | Snorlax blocks `NORTH` until defeated |
| `crossroad` | Lavender Crossroad | Leftovers item |
| `tower` | Pokemon Tower | Cubone's Skull |
| `lavender_town` | Lavender Town | Mr. Fuji (quest-giver) |

### Items

| Item | Type | Effect |
|------|------|--------|
| Thick Bone | Weapon | Sets player damage to `34` |
| Leftovers | Quest item | Required for `quest_blocked_path` |
| Cubone's Skull | Quest item | Required for `quest_tower_mystery` |
| Poke Flute | Quest reward | Granted on main quest completion |

### Quests

| ID | Type | Objective | Reward |
|----|------|-----------|--------|
| `quest_blocked_path` | deliver | Defeat Snorlax -> bring Leftovers to the Club President | Poke Flute |
| `quest_tower_mystery` | collect | Retrieve Cubone's Skull for Mr. Fuji | — |

---

## ⚔️ Combat System

Turn-based, player-initiated combat against hostile NPCs.

```
ATTACK <npc>
      │
      ▼
 [combat state]
      │
      ├─ USE_ITEM  -> player deals weapon dmg (34) or unarmed dmg (5)
      │                └─ NPC counter-attacks with attack_dmg
      ├─ DEFEND    -> same as USE_ITEM but incoming damage is halved
      ├─ FLEE      -> exits combat, no damage taken
      └─ STATUS    -> shows current HP and combat state
```

- Players start at **100 HP**. Reaching 0 HP triggers respawn at `start` with `Max_HP - 1` (99) HP.
- Combat is **fully deterministic** — no random component.
- All events are pushed as `EVT COMBAT ...`; victories are broadcast to the room.

---

## 📋 Requirements

- **Go >= 1.18** (developed and tested with 1.24)
- For the **GUI client** only: a C compiler and X11/OpenGL development libraries (Fyne uses cgo)

---

## 🔍 Error Reference

All error codes are declared in `src/speakserver/errors.go`.
Format: `ERR <code> <SYMBOL>` — first digit identifies the domain.

| Range | Domain |
|-------|--------|
| 1xx | Protocol / command syntax |
| 2xx | Session / authentication |
| 3xx | World / movement |
| 4xx | Items |
| 5xx | NPCs |
| 6xx | Combat |
| 7xx | Quests |
| 8xx | Groups |
| 9xx | Server |

<details>
<summary>Full error code table</summary>

| Code | Symbol | Meaning |
|------|--------|---------|
| 100 | `MALFORMED_COMMAND` | The command could not be parsed |
| 101 | `UNKNOWN_COMMAND` | Command not recognized |
| 102 | `MISSING_ARGUMENT` | A required argument is missing |
| 103 | `UNEXPECTED_ARGUMENT` | An argument was given to a command that takes none |
| 104 | `INVALID_ARGUMENT` | The argument value is not valid |
| 105 | `MESSAGE_TOO_LONG` | Chat message exceeds the allowed length |
| 106 | `CONTROL_D` | Connection closed via Ctrl-D |
| 200 | `NAME_IN_USE` | The requested name is already connected |
| 201 | `NAME_TOO_SHORT` | Name has fewer than 3 characters |
| 202 | `NAME_TOO_LONG` | Name has more than 12 characters |
| 203 | `NAME_INVALID` | Name contains non-letter characters |
| 204 | `CONNECTION_TIMEOUT` | Connection timed out (auth or inactivity) |
| 300 | `NO_EXIT` | No exit in that direction |
| 301 | `NOT_IN_ROOM` | Player is not currently in a room |
| 302 | `PATH_BLOCKED` | An NPC blocks that exit |
| 400 | `ITEM_NOT_FOUND` | No matching item in room/inventory |
| 401 | `ITEM_NOT_OBTAINABLE` | The item cannot be taken |
| 402 | `HANDS_FULL` | Player already holds a weapon |
| 500 | `NPC_NOT_FOUND` | No matching NPC in the room |
| 501 | `NPC_NO_DIALOGUE` | The NPC has no dialogue |
| 502 | `NPC_NOT_HOSTILE` | Cannot attack a non-hostile NPC |
| 503 | `NPC_HOSTILE` | Cannot talk to a hostile NPC |
| 600 | `NOT_IN_COMBAT` | Action requires being in combat |
| 601 | `ALREADY_IN_COMBAT` | NPC is already engaged in combat |
| 602 | `TARGET_GONE` | Combat target no longer exists |
| 603 | `TARGET_DEFEATED` | NPC is already defeated |
| 604 | `COMMAND_NOT_ALLOWED_IN_COMBAT` | Only `USE_ITEM`/`DEFEND`/`FLEE`/`STATUS` allowed in combat |
| 700 | `QUEST_NOT_FOUND` | No quest with that ID |
| 701 | `QUEST_ALREADY_ACTIVE` | Quest is already in progress |
| 702 | `QUEST_ALREADY_COMPLETED` | Quest is already completed |
| 703 | `QUEST_NOT_ACTIVE` | Quest is not in progress |
| 704 | `OBJECTIVE_INCOMPLETE` | Quest objective not yet met |
| 705 | `MISSING_REQUIRED_ITEM` | Required item not in inventory |
| 800 | `NOT_IN_GROUP` | Player is not in a group |
| 801 | `ALREADY_IN_GROUP` | Player is already in a group |
| 802 | `GROUP_NOT_FOUND` | No group matches the request |
| 803 | `NOT_GROUP_LEADER` | Only the group leader can perform this action |
| 804 | `USER_NOT_FOUND` | No connected user with that name |
| 900 | `INTERNAL_ERROR` | Unexpected server-side error |

</details>

---

## 📊 Server Logging

The server emits **structured JSON logs** to stdout via `src/logger`.

```json
{
  "time": "2026-06-22T18:04:11.512Z",
  "level": "INFO",
  "msg": "command received",
  "player": "alice",
  "addr": "127.0.0.1:53024",
  "cmd": "TAKE",
  "args": "item_thick_bone"
}
```

Filter for warnings and errors only:

```bash
./bin/tap-server | jq 'select(.level=="WARN")'
```

Abuse detection monitors **command flooding** (> 20 commands / 10 s per connection) and **rapid connections** (> 5 connections / 10 s per IP).

---

## 📁 Project Structure

```
TAP/
├── cmd/
│   ├── server/         # TCP server entry point
│   ├── client/         # CLI client entry point
│   └── client-gui/     # GUI client entry point
├── src/
│   ├── game/           # Command handlers
│   ├── models/         # Player, Room, NPC, Item, World structs
│   ├── network/        # ClientAtender, Hub, Authentication
│   ├── parse/          # Command parser
│   ├── speakserver/    # OK/ERR/EVT reply format and error catalog
│   └── logger/         # Structured JSON logger
├── assets/             # World workflow diagram
├── data.yaml           # World definition (rooms, items, NPCs, quests)
├── Makefile
└── go.mod
```

---

## 🛠️ Tech Stack

| Component | Technology |
|-----------|-----------|
| Language | Go 1.18+ |
| GUI framework | [Fyne](https://fyne.io/) |
| Concurrency | Goroutines + channels + `sync.RWMutex` |
| World definition | YAML (`data.yaml`) |
| Logging | `log` + `encoding/json` (stdlib) |
| Build tool | GNU Make |

---

## 👥 Authors

| Contributor | Scope |
|-------------|-------|
| **dmena-li** | TCP server foundation, `Hub`/`ClientAtender` lifecycle, protocol framing, world model, command parser, gameplay (movement, combat, quests, groups), initial Fyne GUI |
| **egalindo** | Build tooling, lint/formatting, error-code catalog, world redesign (Kanto arc + Snorlax mechanic), structured JSON logging, robustness fixes, GUI V4 polish, documentation |

---

## 📄 License

This project is open source. See [LICENSE](LICENSE) for details.
