package role

type ID int64

type Type string

const (
	Admin  Type = "admin"
	Vendor Type = "vendor"
	User   Type = "user"
	Guest  Type = "guest"
)
