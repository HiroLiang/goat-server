package shared

type RequestContext struct {
	IP      string
	TraceID string
}

type AuthContext struct {
	UserID string
	RoleID string
	Token  string
}

type BaseInput struct {
	Request RequestContext
	Auth    *AuthContext
}
type UseCaseInput[T any] struct {
	Base BaseInput
	Data T
}
