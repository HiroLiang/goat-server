package participant

type ID int64

type ParticipantType string

const (
	UserType   ParticipantType = "USER"
	AgentType  ParticipantType = "AGENT"
	SystemType ParticipantType = "SYSTEM"
)
