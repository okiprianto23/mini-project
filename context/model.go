package context

type AuthAccessTokenModel struct {
	RedisAuthAccessTokenModel
	ClientID string `json:"cid"`
}

func (input AuthAccessTokenModel) ConvertToRedisModel() RedisAuthAccessTokenModel {
	return RedisAuthAccessTokenModel{
		ResourceUserID: input.ResourceUserID,
		SignatureKey:   input.SignatureKey,
		Locale:         input.Locale,
		AliasName:      input.AliasName,
		ClientAlias:    input.ClientAlias,
	}
}

type RedisAuthAccessTokenModel struct {
	ResourceUserID int64  `json:"rid"`
	SignatureKey   string `json:"sign"`
	Locale         string `json:"locale"`
	AliasName      string `json:"als"`
	ClientAlias    string `json:"cls"`
}
