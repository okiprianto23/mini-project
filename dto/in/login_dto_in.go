package in

type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`

	Token string `json:"token" empty:"allowed"`
}
