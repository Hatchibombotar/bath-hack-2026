package main

type GenericAction struct {
	Action string `json:"action"`
}

type MessageFriendAction struct {
	Action  string `json:"action"`
	Message string `json:"message"`
}

type PlayerData struct {
	Action   string `json:"action"`
	PlayerId int    `json:"player"`
}

type VisiblePlayerDataAction struct {
	Action     string             `json:"action"`
	PlayerId   int                `json:"player"`
	PlayerData *VisiblePlayerData `json:"visible_player_data"`
}
