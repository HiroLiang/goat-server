package chat

type GetMyGroupsInput struct{}

type GetGroupMessagesInput struct {
	GroupID  int64
	BeforeID *int64
	Limit    uint64
}
