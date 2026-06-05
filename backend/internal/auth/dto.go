package auth

type RegisterRequest struct {
	Username string `json:"username" binding:"required,min=3,max=50,alphanum"`
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=8"`
	FullName string `json:"full_name" binding:"required,min=1,max=100"`
}

type LoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

type GoogleAuthRequest struct {
	GoogleToken string `json:"google_token" binding:"required"`
}

type RefreshRequest struct {
	// refresh token comes from httpOnly cookie, not body
}

type TokenPair struct {
	AccessToken string `json:"access_token"`
}

type AuthResponse struct {
	AccessToken string      `json:"access_token"`
	User        interface{} `json:"user"`
}
