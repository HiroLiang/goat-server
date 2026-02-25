package chatmember

type ID int64

type Role string

const (
	Owner  Role = "OWNER"
	Admin  Role = "ADMIN"
	Member Role = "MEMBER"
	Guest  Role = "GUEST"
)
