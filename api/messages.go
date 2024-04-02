package api

import "fmt"

const (
	REGISTERED_MESSAGE   = `{"type": "registered"}`
	START_GAME_MESSAGE   = `{"type": "game_start"}`
	ENDED_GAME_MESSAGE   = `{"type": "game_ended"}`
	START_TURN_MESSAGE   = `{"type": "turn_start"}`
	ENDED_TURN_MESSAGE   = `{"type": "turn_ended"}`
	VALID_MOVE_MESSAGE   = `{"type": "valid_move"}`
	INVALID_MOVE_MESSAGE = `{"type": "invalid_move"}`
	READING_MOVE_MESSAGE = `{"type": "reading_move"}`

	WIN_MESSAGE  = `{"type": "win"}`
	LOSE_MESSAGE = `{"type": "lose"}`
)

func ErrorMessage(err string) string {
	return fmt.Sprintf(`{"type": "error", "error": "%s"}`, err)
}

func IsValidInitMessage(msg map[string]interface{}) bool {
	typ, ok_typ := msg["type"]
	_, ok_key := msg["key"]

	if !ok_typ || !ok_key {
		return false
	}

	return typ.(string) == "register"
}

func CreateStateMessage(state_information map[string]interface{}, time float64) map[string]interface{} {
	return map[string]interface{}{"type": "state", "state": state_information, "timer": time}
}
