package out

type LoginResponse struct {
	UserID    int64  `json:"userID"`
	Token     string `json:"token"`
	Locale    string `json:"locale"`
	Username  string `json:"username"`
	AliasName string `json:"aliasName"`
}
