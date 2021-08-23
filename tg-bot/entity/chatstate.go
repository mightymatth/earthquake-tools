package entity

type ChatState struct {
	ChatID       int64
	AwaitInput   string
	DisableInput bool
}

type ChatStateUpdate struct {
	AwaitInput   string
	DisableInput bool
}
