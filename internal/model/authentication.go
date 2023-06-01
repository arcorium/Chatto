package model

type TokenType string

const (
	Bearer TokenType = "Bearer"
)

type SignInInput struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type SignUpInput struct {
	Username string `json:"username" binding:"required"` // binding used by gin
	Password string `json:"password" binding:"required"`
}

type TokenDetails struct {
	Id         string `json:"id" gorm:"primarykey"`
	UserId     string `json:"user_id"`
	DeviceName string `json:"device_name"`
	Token      string `json:"token"`
}

type AccessToken struct {
	Type  TokenType `json:"type,omitempty"`
	Token string    `json:"access_token" binding:"required"`
}

type AccessTokenClaims struct {
	UserId    string
	RefreshId string
	//Exp       uint64
}
