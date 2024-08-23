package user

import (
	"main-xyz/context"
	"main-xyz/dto/out"
	"main-xyz/repository"
	"main-xyz/router"
)

func (u userService) GetListUser(ctx *context.ContextModel, param router.URLParam, dto interface{}) (headers map[string]string, output interface{}, err error) {
	var dbResult []interface{}
	// get list query user
	dbResult, err = u.userDAO.GetListUser(ctx)
	if err != nil {
		return
	}

	var result []out.UserResponse
	for _, d := range dbResult {
		temp := d.(repository.UserInformation)
		result = append(result, out.UserResponse{
			ID:             temp.ID.Int64,
			AliasName:      temp.AliasName.String,
			Username:       temp.Username.String,
			Locale:         temp.Locale.String,
			ResourceUserID: temp.ResourceUserID.Int64,
			ClientID:       temp.ClientID.String,
			Email:          temp.Email.String,
			AuthUserID:     temp.AuthUserID.Int64,
		})
	}

	output = out.DefaultResponsePayloadMessage{
		Status: out.DefaultResponsePayloadStatus{
			Code:    "OK",
			Message: u.bundles.ReadMessageBundle("user", "SUCCESS_LIST_MESSAGE", ctx.AuthAccessTokenModel.Locale, nil),
		},
		Data: result,
	}

	return
}
