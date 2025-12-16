package user

// RegisterInput represents the payload required for user registration.
type RegisterInput struct {
	Name     string
	Email    string
	Password string
}

// LoginInput represents the required fields for a user login request.
// It includes the user's email and password for authentication.
type LoginInput struct {
	Email    string
	Password string
}
