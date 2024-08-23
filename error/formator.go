package error

import (
	"main-xyz/constanta"
	"main-xyz/context"
	"main-xyz/dto/out"
	"main-xyz/error/bundles"
	"strings"
	"time"
)

var defaultLanguage = "en-US"
var appVersion = "1.0.0"
var defaultInternalCode = "E-5-CMD-SRV-001"
var _formator *formator

func NewErrorFormator(
	bundles bundles.Bundles,
) Formator {
	_formator = &formator{
		bundles: bundles,
	}
	return _formator
}

type errorMessageParam struct {
	err error
	ctx *context.ContextModel
}

func (e *errorMessageParam) WithContext(
	ctx *context.ContextModel,
) *errorMessageParam {
	e.ctx = ctx
	return e
}

type formator struct {
	bundles bundles.Bundles
}

func (f *formator) DefaultLanguage(
	lang string,
) *formator {
	defaultLanguage = lang
	return f
}

func (f *formator) Version(
	version string,
) *formator {
	appVersion = version
	return f
}

func (f *formator) DefaultInternalCode(
	defaultCode string,
) *formator {
	defaultInternalCode = defaultCode
	return f
}

func NewErrorMessageParam(
	err error,
) *errorMessageParam {
	return &errorMessageParam{
		err: err,
	}
}

func (f formator) ReformatErrorMessage(
	param errorMessageParam,
) out.DefaultErrorResponse {

	result := out.DefaultErrorResponse{
		DefaultMessage: out.DefaultMessage{
			Success: false,
		},
	}

	language := defaultLanguage

	result.DefaultMessage.Header = out.Header{
		Version:   appVersion,
		Timestamp: time.Now().Format(constanta.DefaultDtoOutTimeFormat),
	}

	if param.ctx != nil {
		language = param.ctx.AuthAccessTokenModel.Locale
		result.DefaultMessage.Header.RequestID = param.ctx.ClientAccess.Logger.ModelLogger.RequestID
	}

	switch errs := param.err.(type) {
	case *UnbundledErrorMessages:
		result.Payload = out.DefaultError{
			Status: errs.status,
			Code:   errs.code.Error(),
		}

		param := make(map[string]interface{})

		if errs.reason != "" {
			result.Payload.Message = errs.reason
			return result
		}

		if errs.function != nil {
			tempParam := errs.function(errs.param...)

			for key := range tempParam {
				param[key] = tempParam[key].Param
				val, _ := tempParam[key].Param.(string)

				if tempParam[key].IsConverted {
					param[key] = f.bundles.ReadMessageBundle("common.constanta", strings.ToUpper(val), language, nil)
				}
			}
		}

		result.Payload.Message = f.bundles.ReadMessageBundle("common.error", errs.Error(), language, param)
	default:
		result.Payload = out.DefaultError{
			Status:  500,
			Code:    defaultInternalCode,
			Message: f.bundles.ReadMessageBundle("common.error", defaultInternalCode, language, nil),
		}
	}

	return result
}
