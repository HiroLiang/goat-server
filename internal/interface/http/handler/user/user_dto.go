package user

// RegisterRequest represents the payload required for user registration.
type RegisterRequest struct {
	Name     string `json:"name" binding:"required"`
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6"`
}

// RegisterResponse User register response
type RegisterResponse struct {
	Message string `json:"message"`
}

// LoginRequest represents the required fields for a user login request.
// It includes the user's email and password for authentication.
type LoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6"`
}

// LoginResponse represents the server's response after a successful login.
// It contains a token and user details.
type LoginResponse struct {
	Message string `json:"message"`
}

// CurrentUserResponse queried for user login request.
type CurrentUserResponse struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	CreateAt string `json:"create_at"`
}
