package entity

type ChatState struct {
	ChatID int64
	State  string
}

type ChatStateUpdate struct {
	State  string
}
