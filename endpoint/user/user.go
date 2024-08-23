package user

import (
	"main-xyz/router"
	"main-xyz/service/user"
	"net/http"
)

func NewUserEndpoint(
	httpValidator *router.HTTPController,
	srv user.UserService,
) *userEndpoint {
	return &userEndpoint{
		httpValidator: httpValidator,
		srv:           srv,
	}
}

type userEndpoint struct {
	httpValidator *router.HTTPController
	srv           user.UserService
}

func (e *userEndpoint) RegisterEndpoint() {

	e.httpValidator.HandleFunc(
		router.NewHandleFuncParam(
			"/main/user",
			e.httpValidator.WrapService(
				router.NewWarpServiceParam(
					e.srv,
					e.srv.GetListUser,
					e.httpValidator.ControllerValidator.UserAccessValidator,
				),
			),
			http.MethodGet,
		),
	)
}
