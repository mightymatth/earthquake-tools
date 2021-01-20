package entity

type ChatState struct {
	ChatID     int64
	AwaitInput string
}

type ChatStateUpdate struct {
	AwaitInput string
}
