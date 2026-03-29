package main

type GenericAction struct {
	Action string `json:"action"`
}

type MessageFriendAction struct {
	Action  string `json:"action"`
	Message string `json:"message"`
}

type DuckCustomisationAction struct {
	Action       string `json:"action"`
	DuckName     string `json:"duckname"`
	DuckSkinFile string `json:"file"`
	Player       string `json:"playername"`
}

type WorkingStatus struct {
	Action    string `json:"action"`
	Player    string `json:"playername"`
	IsWorking bool   `json:"isworking"`
}
