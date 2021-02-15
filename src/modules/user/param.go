package user

// NewUser is parameter when creating new user
type NewUser struct {
	Email    string `json:"email"`
	Password string `json:"password"`
	Fullname string `json:"fullname"`
}
