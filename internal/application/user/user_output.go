package user

// LoginOutput represents the server's response after a successful login.
// It contains a token and user details.
type LoginOutput struct {
	Token string
}

// CurrentUserOutput let current user logout
type CurrentUserOutput struct {
	Name     string
	Email    string
	CreateAt string
}
