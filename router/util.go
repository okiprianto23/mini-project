package router

import (
	"context"
	"encoding/json"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"main-xyz/constanta"
	internalCtx "main-xyz/context"
	error2 "main-xyz/error"
	"net/http"
	"net/url"
	"strings"
)

func (g HTTPController) preStart(
	r *http.Request,
	param *wrapServiceParam,
) error {

	if param.multipart {
		return g.reader.ParseMultipartForm(r)
	}

	return nil
}

func (g HTTPController) getContext(
	r *http.Request,
) *internalCtx.ContextModel {

	var ctx *internalCtx.ContextModel

	rCtx := r.Context().Value(constanta.ApplicationContextConstanta)
	if rCtx == nil {

		ctx = internalCtx.NewContextModel()
		context := context.WithValue(
			r.Context(),
			constanta.ApplicationContextConstanta,
			ctx,
		)

		r = r.WithContext(context)
	} else {
		ctx = rCtx.(*internalCtx.ContextModel)
	}

	requestID := ReadHeader(r, constanta.RequestIDConstanta)
	if requestID == "" {
		requestID = GetUUID()
		r.Header.Set(constanta.RequestIDConstanta, requestID)
	}

	ctx.ClientAccess.Logger.Set("request_id", requestID)
	ctx.ClientAccess.Logger.Set("ip", ReadHeader(r, constanta.IPAddressConstanta))
	ctx.ClientAccess.Logger.Set("source", ReadHeader(r, constanta.SourceConstanta))
	ctx.ClientAccess.Logger.Set("access_token", ReadHeader(r, constanta.AuthorizationHeaderConstanta))

	return ctx

}

func GetUUID() (output string) {
	UUID, _ := uuid.NewRandom()
	output = UUID.String()
	output = strings.Replace(output, "-", "", -1)
	return
}

func (g HTTPController) readBody(
	ctx *internalCtx.ContextModel,
	param *wrapServiceParam,
	request *http.Request,
) (
	dto interface{},
	err error,
) {
	var stringBody string

	if request.Method != "GET" {
		if !param.readBody {
			return
		}

		dto = param.service.GetDTO()

		var byteIn int
		stringBody, byteIn, err = ReadBody(request)
		if err != nil {
			return nil, error2.ErrReadBody
		}

		ctx.ClientAccess.Logger.Set(constanta.LoggerByteIn, byteIn)

		err = json.Unmarshal([]byte(stringBody), dto)
		if err != nil {
			return nil, error2.ErrMarshalingBody
		}

		err = g.validator.ValidateByTag(ctx, dto, param.functionValidator, param.menu)
		if err != nil {
			return nil, err
		}
	}

	return
}

func (g HTTPController) getURLParam(
	request *http.Request,
	path []string,
	header []string,
) (
	URLParam,
	map[string]string,
) {
	param := URLParam{}

	param.Query = g.readQueryParam(request)
	param.Path = g.readPathParam(request, path)

	return param, g.readHeaders(request, header)
}

func (g HTTPController) readQueryParam(
	request *http.Request,
) map[string]string {
	return g.generateQueryParam(request)
}

func (g HTTPController) readPathParam(
	request *http.Request,
	params []string,
) map[string]string {

	path := make(map[string]string)

	for i := 0; i < len(params); i++ {
		path[params[i]] = mux.Vars(request)[params[i]]
	}

	return path
}

func (g HTTPController) readHeaders(
	request *http.Request,
	headers []string,
) map[string]string {

	header := make(map[string]string)

	for i := 0; i < len(headers); i++ {
		header[headers[i]] = request.Header.Get(headers[i])
	}

	return header
}

func (g HTTPController) generateQueryParam(request *http.Request) map[string]string {
	result := make(map[string]string)
	defer func() {
		_ = recover()
	}()

	var errs error
	rawQuery := request.URL.RawQuery
	rawSplit := strings.Split(rawQuery, "&")

	for key := range rawSplit {
		splitEqual := strings.Split(rawSplit[key], "=")
		result[splitEqual[0]], errs = url.QueryUnescape(splitEqual[1])
		if errs != nil {
			result[splitEqual[0]] = splitEqual[1]
		}
	}

	return result
}
