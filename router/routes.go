package router

import (
	"context"
	"github.com/go-redis/redis/v7"
	"github.com/gorilla/mux"
	"main-xyz/constanta"
	error2 "main-xyz/error"
	"main-xyz/server"
	"main-xyz/token"
	"main-xyz/utils/validator/basic_validator"
	"main-xyz/utils/validator/multipart_validator"
	"main-xyz/utils/validator/tag_validator"
	"net/http"
)

type HandleFuncParam struct {
	path   string
	f      func(http.ResponseWriter, *http.Request)
	method []string
}

func NewHandleFuncParam(
	path string,
	f func(http.ResponseWriter, *http.Request),
	method ...string,
) *HandleFuncParam {
	return &HandleFuncParam{
		path:   path,
		f:      f,
		method: method,
	}
}

type HTTPController struct {
	ControllerValidator
	router             *mux.Router
	version            string
	reader             server.MultipartReader
	basicValidator     basic_validator.BasicValidator
	validator          tag_validator.TagValidator
	multipartValidator multipart_validator.MultipartValidator
	formator           error2.Formator
}

func NewHTTPController() *HTTPController {
	validator := &HTTPController{
		version: "1.0.0",
	}

	knownHeader = map[string]string{
		constanta.AuthorizationHeaderConstanta:          "",
		constanta.RequestIDConstanta:                    "",
		constanta.IPAddressConstanta:                    "",
		constanta.SourceConstanta:                       "",
		constanta.TimestampSignatureHeaderNameConstanta: "",
		constanta.SignatureHeaderNameConstanta:          "",
		constanta.DeviceHeaderConstanta:                 "",
		constanta.ResourceHeaderConstanta:               "",
		constanta.RedirectURINameConstanta:              "",
		constanta.RedirectMethodNameConstanta:           "",
		constanta.ClientRequestTimestamp:                "",
	}

	return validator
}

func (g HTTPController) HandleFunc(param *HandleFuncParam) {
	g.router.HandleFunc(param.path, param.f).Methods(param.method...)
}

func (g *HTTPController) Router(router *mux.Router) *HTTPController {
	g.router = router
	return g
}

func (g *HTTPController) Version(version string) *HTTPController {
	g.version = version
	return g
}

func (cv *ControllerValidator) TokenValidator(tokenValidator token.UserJWTValidator) *ControllerValidator {
	cv.tokenValidator = tokenValidator
	return cv
}

func (g *HTTPController) Formator(formator error2.Formator) *HTTPController {
	g.formator = formator
	return g
}

func (g *HTTPController) TagValidator(validator tag_validator.TagValidator) *HTTPController {
	g.validator = validator
	return g
}

func (g *HTTPController) BasicValidator(validator basic_validator.BasicValidator) *HTTPController {
	g.basicValidator = validator
	return g
}

func (g *HTTPController) MultipartValidator(validator multipart_validator.MultipartValidator) *HTTPController {
	g.multipartValidator = validator
	return g
}

func (g *HTTPController) MultipartReader(reader server.MultipartReader) *HTTPController {
	g.reader = reader
	return g
}

func (g *HTTPController) Redis(redis *redis.Client) *HTTPController {
	g.redis = redis
	return g
}

var knownHeader map[string]string

func (g HTTPController) Converter(
	r *http.Request,
) (
	context.Context,
	map[string]string,
) {
	result := make(map[string]string)

	for keys := range knownHeader {
		value := r.Header.Get(keys)
		if value != "" {
			result[keys] = value
		}
	}

	if r.Context() != nil {
		r = r.WithContext(context.Background())
	}

	return r.Context(), result
}
