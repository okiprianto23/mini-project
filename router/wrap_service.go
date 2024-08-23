package router

import (
	"context"
	"encoding/json"
	"github.com/gorilla/mux"
	"main-xyz/constanta"
	internalCtx "main-xyz/context"
	"main-xyz/dto/out"
	error2 "main-xyz/error"
	"main-xyz/service"
	"main-xyz/utils/validator/multipart_validator"
	"net/http"
	"strings"
	"time"
)

type FunctionServe func(
	*internalCtx.ContextModel,
	URLParam,
	interface{},
) (
	map[string]string,
	interface{},
	error,
)

type ServerAccessValidator func(ctx context.Context, header map[string]string) error

type URLParam struct {
	Query map[string]string
	Path  map[string]string
}

type wrapServiceParam struct {
	serve             FunctionServe
	cv                ServerAccessValidator
	service           service.Services
	readBody          bool
	menu              string
	pathParams        []string
	headers           []string
	functionValidator string
	multipart         bool
}

func (w *wrapServiceParam) Multipart() *wrapServiceParam {
	w.multipart = true
	return w
}

func (w *wrapServiceParam) Headers(headers ...string) *wrapServiceParam {
	w.headers = headers
	return w
}

func (w *wrapServiceParam) PathParams(paths ...string) *wrapServiceParam {
	w.pathParams = paths
	return w
}

func (w *wrapServiceParam) NotReadBody() *wrapServiceParam {
	w.readBody = false
	return w
}

func (w *wrapServiceParam) Menu(menu string) *wrapServiceParam {
	w.menu = menu
	return w
}

func (w *wrapServiceParam) FunctionValidator(functionValidator string) *wrapServiceParam {
	w.functionValidator = functionValidator
	return w
}

func NewWarpServiceParam(
	service service.Services,
	serve FunctionServe,
	cv ServerAccessValidator,
) *wrapServiceParam {
	return &wrapServiceParam{
		cv:       cv,
		service:  service,
		serve:    serve,
		readBody: true,
	}
}

func (g *HTTPController) WrapService(
	param *wrapServiceParam,
) func(
	http.ResponseWriter,
	*http.Request,
) {
	return func(rw http.ResponseWriter, r *http.Request) {
		var (
			ctx      *internalCtx.ContextModel
			err      error
			payload  interface{}
			header   map[string]string
			urlParam URLParam
			timeNow  = time.Now()
		)

		defer func() {
			g.response(
				ctx,
				r,
				rw,
				timeNow,
				header,
				payload,
				err,
				urlParam,
			)
		}()

		ctx = g.getContext(r)

		err = g.preStart(
			r,
			param,
		)

		if err != nil {
			err = error2.ErrReadBody
			return
		}

		ctx.ClientAccess.Path = r.URL.Path

		//validate tokn
		_, headers := g.Converter(r)

		ctx.ClientAccess.Timestamp, err = time.Parse(constanta.DefaultDtoOutTimeFormat, headers[constanta.ClientRequestTimestamp])
		if err != nil {
			ctx.ClientAccess.Timestamp = time.Now()
		}

		//check token to redis
		err = param.cv(ctx.ToContext(), headers)
		if err != nil {
			return
		}

		c, valid := r.Context().Value(
			constanta.ApplicationContextConstanta,
		).(*internalCtx.ContextModel)

		if valid {
			ctx = c
		}

		var dto interface{}

		if !param.multipart {
			dto, err = g.readBody(
				ctx,
				param,
				r,
			)
			if err != nil {
				return
			}
		} else {
			dto, err = g.multipartValidator.ValidateMultipart(
				multipart_validator.NewValidateMultipartParam(
					ctx,
					r,
					param.service.GetMultipartDTO(),
				).Menu(
					param.menu,
				).FunctionValidator(
					param.functionValidator,
				),
			)
		}

		//get url param
		urlParam, ctx.ClientAccess.Headers = g.getURLParam(r, param.pathParams, param.headers)

		header, payload, err = param.serve(
			ctx,
			urlParam,
			dto,
		)

	}
}

func (g HTTPController) response(
	ctx *internalCtx.ContextModel,
	r *http.Request,
	rw http.ResponseWriter,
	timeNow time.Time,
	headers map[string]string,
	payload interface{},
	err error,
	pathParam URLParam) {
	var errs out.DefaultErrorResponse
	//var usedErr = err
	statusCode := 200

	defer func() {
		usedPath, _ := mux.CurrentRoute(r).GetPathTemplate()
		if usedPath == "" {
			usedPath = r.URL.Path
		}
		ctx.ClientAccess.Logger.Set(constanta.LoggerProcessingTime, time.Since(timeNow).Microseconds())
		ctx.ClientAccess.Logger.Set(constanta.LoggerUrl, usedPath)
		ctx.ClientAccess.Logger.Set(constanta.LoggerMethod, r.Method)

		msg := "Api Called"
		if errs.Payload.Message != "" {
			msg = errs.Payload.Message
			ctx.ClientAccess.Logger.Set(constanta.LoggerCode, errs.Payload.Code)
		}

		// comment agar tidak 2 kali logger
		//if usedErr != nil {
		//	ctx.ClientAccess.Logger.Logger.Error(usedErr.Error())
		//}

		ctx.ClientAccess.Logger.Logger.Info(msg)
	}()

	rw.Header().Set(constanta.RequestIDConstanta, ctx.ClientAccess.Logger.ModelLogger.RequestID)

	for key := range headers {
		headersList := rw.Header().Get("Access-Control-Allow-Headers") + ", " + strings.ToLower(key)
		rw.Header().Set("Access-Control-Allow-Headers", headersList)
		rw.Header().Add(key, headers[key])
	}

	var data []byte
	var length int
	if err != nil {
		errs = g.formator.ReformatErrorMessage(
			*error2.NewErrorMessageParam(
				err,
			).WithContext(
				ctx,
			),
		)

		statusCode = errs.Payload.Status
		payload = errs
	} else {
		result := out.DefaultResponse{
			DefaultMessage: out.DefaultMessage{
				Success: true,
			},
		}
		result.DefaultMessage.Header = out.Header{
			Version:   g.version,
			Timestamp: time.Now().Format(constanta.DefaultDtoOutTimeFormat),
		}

		if ctx != nil {
			result.DefaultMessage.Header.RequestID = ctx.ClientAccess.Logger.ModelLogger.RequestID
		}

		result.Payload = payload
		payload = result
	}

	rw.Header().Set("Content-Type", "application/json")
	ctx.ClientAccess.Logger.Set(constanta.LoggerStatus, statusCode)
	rw.WriteHeader(statusCode)

	if data == nil {
		data, err = json.Marshal(payload)
		if err != nil {
			ctx.ClientAccess.Logger.Logger.Error("Error on Marshaling Message")
			return
		}

		//if statusCode == http.StatusOK {
		//	timestamp := time.Now().Format(constanta.DefaultTimeFormat)
		//	digest := GenerateMessageDigest(string(data))
		//	signature := GenerateSignature(
		//		r.Method,
		//		r.URL.Path,
		//		ctx.ClientAccess.Logger.ModelLogger.AccessToken,
		//		digest,
		//		timestamp,
		//		ctx.AuthAccessTokenModel.SignatureKey,
		//	)
		//
		//	rw.Header().Set(constanta.TimestampSignatureHeaderNameConstanta, timestamp)
		//	rw.Header().Set(constanta.SignatureHeaderNameConstanta, signature)
		//}

		length, err = rw.Write(data)
		if err != nil {
			ctx.ClientAccess.Logger.Logger.Error("Error on Writing Message")
			return
		}
	}
	ctx.ClientAccess.Logger.Set(constanta.LoggerByteOut, length)
}
