package token

import "github.com/golang-jwt/jwt/v4"

type PayloadJWTToken struct {
	ClientID string `json:"cid"`
	Resource string `json:"resource"`
	Scope    string `json:"scope"`
	Locale   string `json:"locale"`
	jwt.RegisteredClaims
}
