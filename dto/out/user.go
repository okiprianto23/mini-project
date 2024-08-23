package out

type UserResponse struct {
	ID             int64  `json:"ID"`
	AliasName      string `json:"aliasName"`
	Username       string `json:"username"`
	Locale         string `json:"locale"`
	ResourceUserID int64  `json:"resourceUserID"`
	ClientID       string `json:"clientID"`
	Email          string `json:"email"`
	AuthUserID     int64  `json:"authUserID"`
}
