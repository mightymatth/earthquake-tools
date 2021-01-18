package entity

type ChatState struct {
	ChatID     int64
	AwaitInput AwaitInput
}

type ChatStateUpdate struct {
	AwaitInput AwaitInput
}

type AwaitInput string

const (
	CreateSubName AwaitInput = "CREATE_SUB_NAME"
)
