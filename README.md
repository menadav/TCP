# TAP — The Answer Protocol

> A real-time multiplayer text adventure engine built on a custom TCP protocol, written in Go.

[![Go](https://img.shields.io/badge/Go-1.18%2B-00ADD8?style=flat&logo=go)](https://go.dev/)
[![License](https://img.shields.io/badge/license-MIT-blue.svg)](LICENSE)
[![Build](https://img.shields.io/badge/build-make-success)](#building)

---

TAP is a shared-world, retro text adventure (MUD) where multiple players connect simultaneously to a persistent server and explore rooms, chat, form groups, fight NPCs and complete quests — all in real time over a line-based TCP protocol.

The project ships three components: a **TCP server**, a **CLI client** and a **GUI client** (built with [Fyne](https://fyne.io/)).

---

## Table of Contents

- [Features](#features)
- [Architecture](#architecture)
- [Protocol](#protocol)
- [World Design](#world-design)
- [Getting Started](#getting-started)
- [Make Targets](#make-targets)
- [Error Reference](#error-reference)
- [Server Logging](#server-logging)
- [Testing](#testing)
- [Authors](#authors)

---

## Features

- **Multiplayer real-time** — unlimited concurrent TCP connections; all state changes are broadcast live.
- **Dual clients** — lightweight CLI and a Fyne-based GUI with live player counters, separated chat/log views and D-pad movement.
- **Custom protocol** — line-based `OK` / `ERR` / `EVT` framing with a versioned handshake (`proto=1`).
- **Turn-based combat** — deterministic damage, defend/flee mechanics, NPC counter-attacks, and auto-respawn.
- **Quest system** — `collect`, `kill`, `explore` and `deliver` quest types, tracked per player.
- **Group system** — create, invite, join and leave named groups; group-scoped chat and events.
- **Structured JSON logging** — RFC3339 timestamps, levelled entries, abuse detection (flood + rapid connections).
- **Data-driven world** — rooms, items, NPCs and quests defined in `data.yaml`; no hardcoded world state.

---

## Architecture

The server uses a **central hub + one goroutine per client** model. All global state mutations are serialised through a single `Hub` goroutine; per-entity `sync.RWMutex` locks protect `Room`, `Player` and `Group` structs.

```
cmd/server ──► network.ClientAtender        (1 goroutine / client)
                     │
                     ├─ network.Authentication   (CONNECT handshake)
                     ├─ parse.ParseCommandCli ──► src/game/*  (handlers)
                     └─ player.ListenMsg          (1 goroutine / client → EVT)
                     ▲
   models.Hub.Run ───┘   (single goroutine: Register / Unregister / Broadcast)
```

Key design decisions:

- **No shared mutable global state** — the client map lives exclusively inside `Hub.Run`; connections talk to it via typed channels (`Register`, `Unregister`, `Broadcast`).
- **Buffered per-player `MsgChan`** — asynchronous event delivery keeps clients responsive without blocking command processing.
- **Inline command dispatch** — `parse.ParseCommandCli` is a `switch` over a small, fixed command set; readable and easy to audit.

---

## Protocol

Transport: **TCP, UTF-8, one message per line** (`\n` terminated).

### Handshake

```
← OK hello proto=1
→ CONNECT <name>        # 3–12 letters, unique across the session
← OK welcome name=<name>
```

### Message prefixes

| Prefix | Meaning |
|--------|---------|
| `OK <payload>` | Success — payload is `key=value` or a JSON object |
| `ERR <code> <SYMBOL>` | Failure — see [Error Reference](#error-reference) |
| `EVT <category> <data>` | Async event pushed by the server (chat, presence, combat…) |

### Non-standard extensions

| Command | Behaviour |
|---------|-----------|
| `REQ` | Returns a JSON snapshot (`WorldStateResponse`) of the current room — used by the GUI client to refresh its view |
| `DEFEND` / `FLEE` | Combat-only commands (our design choice, not in the base RFC) |
| `WHO` | Replies `OK who={"room":[names],"server":N}` — live player counters |
| `INVENTORY` | Replies with a JSON array of the player's items |

---

## World Design

The world is themed around the **Kanto region** and ships with 8 rooms, 4 items, 3 NPCs and 2 quests.

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
              +-------------+----+   E   +----------------+
              |   silence_bridge |       |  guard_house   |  ← Snorlax blocks NORTH
              | (Silence Bridge) |       | (Guard House)  |
              +--------+---------+       +-------+--------+
                   N   |   S                 N   |   S
              +--------+---------+   E   +-------+--------+
              |     route_11     |------>|   graveyard    |
              |   (Route 11)     |<------| (Graveyard)    |
              +--------+---------+   W   +----------------+
                   S   |   N
              +------------------+
              |      start       |
              | (Vermilion City) |
              +------------------+
```

### Rooms

| Room ID | Name | Notable |
|---------|------|---------|
| `start` | Vermilion City | Spawn point; Club President (quest-giver) |
| `graveyard` | Graveyard | Thick Bone weapon |
| `route_11` | Route 11 | — |
| `guard_house` | Guard House | — |
| `silence_bridge` | Silence Bridge | Snorlax blocks `NORTH` until defeated |
| `crossroad` | Lavender Crossroad | Leftovers item |
| `tower` | Pokémon Tower | Cubone's Skull |
| `lavender_town` | Lavender Town | Mr. Fuji (quest-giver) |

### Items

| Item | Type | Effect |
|------|------|--------|
| Thick Bone | Weapon | Sets player damage to `34` |
| Leftovers | Quest item | Required for `quest_blocked_path` |
| Cubone's Skull | Quest item | Required for `quest_tower_mystery` |
| Poké Flute | Quest reward | Granted on main quest completion |

### Quests

| ID | Type | Objective | Reward |
|----|------|-----------|--------|
| `quest_blocked_path` | deliver | Defeat Snorlax → bring Leftovers to the Club President | Poké Flute |
| `quest_tower_mystery` | collect | Retrieve Cubone's Skull for Mr. Fuji | — |

### Combat

- Players start at **100 HP**. `ATTACK <npc>` initiates combat and sets the player to `combat` state.
- Allowed commands while in combat: `USE_ITEM` (attack), `DEFEND`, `FLEE`, `STATUS`.
- **Damage:** fully deterministic — `5` unarmed, or the equipped weapon's `dmg` value.
- **Counter-attack:** the NPC strikes back after each non-lethal hit (Snorlax deals `49`). `DEFEND` halves incoming damage for that turn.
- **Defeat:** reaching 0 HP respawns the player at `start` with `Max_HP − 1` (99) HP.
- All combat events are pushed as `EVT COMBAT ...`; victories are broadcast to the room.

---

## Getting Started

### Requirements

- **Go ≥ 1.18** (developed and tested with 1.24)
- For the **GUI client** only: a C compiler and X11/OpenGL development libraries (Fyne uses cgo)

### Quick start

```bash
# 1. Fetch dependencies
make install

# 2. Build all binaries into ./bin
make build

# 3. Start the server (terminal 1)
make run-server

# 4. Connect a CLI client (terminal 2)
make run-client

# 5. Connect a GUI client (terminal 3)
make run-client-gui
```

---

## Make Targets

| Target | Description |
|--------|-------------|
| `make install` | Download module dependencies (`go mod download`) |
| `make build` | Compile server, CLI and GUI into `./bin` |
| `make run-server` | Start the TCP server on `:8080` |
| `make run-client` | Start the CLI client |
| `make run-client-gui` | Start the GUI client |
| `make lint` | `gofmt` check + `go vet` |
| `make fmt-fix` | Reformat source with `gofmt -w` |
| `make clean` | Remove built binaries |

---

## Error Reference

All error codes are declared in `src/speakserver/errors.go` (single source of truth).
Format: `ERR <code> <SYMBOL>`

The leading digit identifies the domain:

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

## Server Logging

The server emits **structured JSON logs** to stdout via `src/logger` (built on `log` + `encoding/json` for Go 1.18 compatibility).

### Log entry format

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

### Event types

| `msg` | Level | Trigger |
|-------|-------|---------|
| `server ready` | INFO | Server starts listening |
| `connection open` | INFO | TCP connection accepted |
| `client registered` | INFO | Successful authentication |
| `auth failed` | WARN | Authentication aborted/failed |
| `connection close` | INFO | Client disconnects |
| `command received` | INFO | Any command from a client |
| `response sent` | INFO | Server reply or event |
| `error response` | WARN | Error code sent to a client |
| `world change` | INFO | Item, NPC or combat state mutation |
| `quest progress` | INFO | Quest accept or complete |
| `abuse detected` | WARN | Command flooding or rapid connections |

### Abuse detection

- **Command flooding** — more than 20 commands within a 10-second window per connection.
- **Rapid connections** — more than 5 connections from the same IP within 10 seconds.

Filter for warnings only:

```bash
./bin/tap-server | jq 'select(.level=="WARN")'
```

---

## Testing

Multiplayer behaviour is validated manually by running the server with multiple concurrent clients.

| Scenario | Steps |
|----------|-------|
| **Presence & chat** | Connect two clients; `MOVE` one between rooms and confirm `EVT ROOM PRESENCE ENTER/LEAVE` on the other. Send `CHAT GLOBAL/ROOM/GROUP` and verify scope-correct delivery. |
| **Items** | `TAKE` an item on client A; confirm it disappears from `LOOK` on client B. `DROP` it and confirm it reappears. |
| **Combat** | `ATTACK` a hostile NPC; cycle through `USE_ITEM`, `DEFEND`, `FLEE`; check `STATUS`; verify respawn at 0 HP. |
| **Quests** | `TALK` to a quest-giver → `QUEST ACCEPT` → fulfil objective → `QUEST COMPLETE`; verify reward in `INVENTORY`. |

---

## Authors

| Contributor | Scope |
|-------------|-------|
| **dmena-li** | Core implementation — TCP server foundation, `Hub`/`ClientAtender` lifecycle, protocol framing, world model, command parser, gameplay (move/look/inventory/combat/quests/groups), initial GUI and Fyne client wiring |
| **egalindo** | Build tooling, lint/formatting, error-code catalog, world redesign (Kanto arc, Snorlax mechanic), structured JSON logging, robustness fixes (hub broadcast, channel-close race), GUI V4 polish, documentation |
