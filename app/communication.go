package main

import "encoding/json"

func (g *Game) SendMessage(message []byte) {
	if g.connectionFail {
		return
	}
	g.sendCh <- message
}

// Runs before update
func (g *Game) HandleMessage(message []byte) error {
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

		g.sendCh <- message
	}

	return nil
}
