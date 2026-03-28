package main

import (
	"encoding/json"

	"github.com/gorilla/websocket"
)

type GenericAction struct {
	Action string `json:"action"`
}

type MessageFriendAction struct {
	Action  string `json:"action"`
	Message string `json:"message"`
}

func handleMessage(player *Player, message []byte) error {
	action := &GenericAction{}
	err := json.Unmarshal(message, action)
	if err != nil {
		return err
	}

	actionName := action.Action

	if actionName == "message_friend" {
		action := &MessageFriendAction{}
		err := json.Unmarshal(message, action)
		if err != nil {
			return err
		}

		message, err := json.Marshal(&GenericAction{Action: "blah blah"})
		if err != nil {
			return err
		}

		err = player.Conn.WriteMessage(websocket.TextMessage, message)
		if err != nil {
			return err
		}

	}

	return nil
}
