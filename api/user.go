package api

// Holds all user related api request and response structs

// User Login

type UserLoginRequest struct {
	AccessToken  string `json:"access_token"`
	IDToken      string `json:"id_token"`
	ExpiresIn    int    `json:"expires_in"`
	TokenType    string `json:"token_type"`
	Scope        string `json:"scope"`
	RefreshToken string `json:"refresh_token"`
}

type UserLoginResp struct{}

// User Get

type UserGetResp struct {
	ID       uint   `json:"id"`
	GoogleID string `json:"google_id"`
	Email    string `json:"email"`
	Name     string `json:"name"`
	Picture  string `json:"picture"`
	Status   string `json:"status"`
	// Add other user data as needed
}

// User Logout

type UserLogoutResp struct{}
