package auth

// LoginParam receive email and password parameters
type LoginParam struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}
