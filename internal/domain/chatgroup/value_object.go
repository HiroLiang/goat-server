package chatgroup

type ID int64

type GroupType string

const (
	Direct  GroupType = "DIRECT"
	Group   GroupType = "GROUP"
	Channel GroupType = "CHANNEL"
	Bot     GroupType = "BOT"
)
