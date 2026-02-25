package chatmessage

type ID int64

type MessageType string

const (
	Text   MessageType = "TEXT"
	Image  MessageType = "IMAGE"
	File   MessageType = "FILE"
	System MessageType = "SYSTEM"
)
