package api

import (
	"encoding/json"
)

func ParseMove(msg []byte) (map[string]interface{}, error) {
	var move map[string]interface{} = nil

	err := json.Unmarshal(msg, &move)

	return move, err
}
