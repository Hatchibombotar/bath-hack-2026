package main

type GenericAction struct {
	Action string `json:"action"`
}

type MessageFriendAction struct {
	Action  string `json:"action"`
	Message string `json:"message"`
}
