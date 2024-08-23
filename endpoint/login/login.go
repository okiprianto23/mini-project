package login

import (
	"main-xyz/router"
	"main-xyz/service/login"
	"net/http"
)

func NewLoginEndpoint(
	httpValidator *router.HTTPController,
	srv login.LoginService,
) *loginEndpoint {
	return &loginEndpoint{
		httpValidator: httpValidator,
		srv:           srv,
	}
}

type loginEndpoint struct {
	httpValidator *router.HTTPController
	srv           login.LoginService
}

func (e *loginEndpoint) RegisterEndpoint() {

	e.httpValidator.HandleFunc(
		router.NewHandleFuncParam(
			"/auth/login",
			e.httpValidator.WrapService(
				router.NewWarpServiceParam(
					e.srv,
					e.srv.LoginService,
					e.httpValidator.ControllerValidator.WhitelistValidator,
				),
			),
			http.MethodPost,
		),
	)

}
