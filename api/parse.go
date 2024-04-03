package api

import (
	"encoding/json"
	"fmt"
)

func ParseMove(msg []byte) map[string]interface{} {
	var move map[string]interface{}

	err := json.Unmarshal(msg, &move)

	if err != nil {
		//Prints the error if not nil
		fmt.Println("Error while decoding the data", err.Error())
		return nil
	}

	return move
}
