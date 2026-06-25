package speak

import (
	"answer_protocol/src/logger"
	"fmt"
	"net"
)

type ErrCode struct {
	Code int
	Sym  string
}

var (
	ErrMalformedCommand   = ErrCode{100, "MALFORMED_COMMAND"}
	ErrUnknownCommand     = ErrCode{101, "UNKNOWN_COMMAND"}
	ErrMissingArgument    = ErrCode{102, "MISSING_ARGUMENT"}
	ErrUnexpectedArgument = ErrCode{103, "UNEXPECTED_ARGUMENT"}
	ErrInvalidArgument    = ErrCode{104, "INVALID_ARGUMENT"}
	ErrMessageTooLong     = ErrCode{105, "MESSAGE_TOO_LONG"}
	ErrCntrlD			  = ErrCode{106, "CONTROL_D"}
	ErrControlChars       = ErrCode{107, "CONTROL_CHARACTERS"}

	ErrNameInUse    = ErrCode{200, "NAME_IN_USE"}
	ErrNameTooShort = ErrCode{201, "NAME_TOO_SHORT"}
	ErrNameTooLong  = ErrCode{202, "NAME_TOO_LONG"}
	ErrNameInvalid  = ErrCode{203, "NAME_INVALID"}
	ErrTimeout      = ErrCode{204, "CONNECTION_TIMEOUT"}

	ErrNoExit      = ErrCode{300, "NO_EXIT"}
	ErrNotInRoom   = ErrCode{301, "NOT_IN_ROOM"}
	ErrPathBlocked = ErrCode{302, "PATH_BLOCKED"}

	ErrItemNotFound  = ErrCode{400, "ITEM_NOT_FOUND"}
	ErrNotObtainable = ErrCode{401, "ITEM_NOT_OBTAINABLE"}
	ErrHandsFull     = ErrCode{402, "HANDS_FULL"}

	ErrNpcNotFound   = ErrCode{500, "NPC_NOT_FOUND"}
	ErrNpcNoDialogue = ErrCode{501, "NPC_NO_DIALOGUE"}
	ErrNpcNotHostile = ErrCode{502, "NPC_NOT_HOSTILE"}
	ErrNpcHostile    = ErrCode{503, "NPC_HOSTILE"}

	ErrNotInCombat     = ErrCode{600, "NOT_IN_COMBAT"}
	ErrAlreadyInCombat = ErrCode{601, "ALREADY_IN_COMBAT"}
	ErrTargetGone      = ErrCode{602, "TARGET_GONE"}
	ErrTargetDefeated  = ErrCode{603, "TARGET_DEFEATED"}
	ErrCommandInCombat = ErrCode{604, "COMMAND_NOT_ALLOWED_IN_COMBAT"}

	ErrQuestNotFound       = ErrCode{700, "QUEST_NOT_FOUND"}
	ErrQuestAlreadyActive  = ErrCode{701, "QUEST_ALREADY_ACTIVE"}
	ErrQuestAlreadyDone    = ErrCode{702, "QUEST_ALREADY_COMPLETED"}
	ErrQuestNotActive      = ErrCode{703, "QUEST_NOT_ACTIVE"}
	ErrObjectiveIncomplete = ErrCode{704, "OBJECTIVE_INCOMPLETE"}
	ErrMissingRequiredItem = ErrCode{705, "MISSING_REQUIRED_ITEM"}

	ErrNotInGroup     = ErrCode{800, "NOT_IN_GROUP"}
	ErrAlreadyInGroup = ErrCode{801, "ALREADY_IN_GROUP"}
	ErrGroupNotFound  = ErrCode{802, "GROUP_NOT_FOUND"}
	ErrNotGroupLeader = ErrCode{803, "NOT_GROUP_LEADER"}
	ErrUserNotFound   = ErrCode{804, "USER_NOT_FOUND"}

	ErrInternal = ErrCode{900, "INTERNAL_ERROR"}
)

func SendErr(conn net.Conn, e ErrCode) {
	fmt.Fprintf(conn, "ERR %d %s\n", e.Code, e.Sym)
	logger.Warn("error response", "addr", logger.Addr(conn), "code", e.Code, "sym", e.Sym)
}
