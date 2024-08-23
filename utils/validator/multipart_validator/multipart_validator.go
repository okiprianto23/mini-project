package multipart_validator

import (
	"encoding/json"
	"main-xyz/constanta"
	error2 "main-xyz/error"
	"main-xyz/server"
	"main-xyz/utils/validator/basic_validator"
	"main-xyz/utils/validator/tag_validator"
	"reflect"
	"strconv"
	"strings"
)

func NewMultipartValidator(
	reader server.MultipartReader,
	basicValidator basic_validator.BasicValidator,
	tagValidator tag_validator.TagValidator,
) MultipartValidator {
	return &multipartValidator{
		reader:         reader,
		basicValidator: basicValidator,
		tagValidator:   tagValidator,
	}
}

type multipartValidator struct {
	reader         server.MultipartReader
	basicValidator basic_validator.BasicValidator
	tagValidator   tag_validator.TagValidator
}

func (m multipartValidator) ValidateMultipart(
	param *validateMultipartParam,
) (
	interface{},
	error,
) {

	reflectType := reflect.TypeOf(param.dto).Elem()
	reflectValue := reflect.ValueOf(param.dto).Elem()

	for i := 0; i < reflectType.NumField(); i++ {
		currentField := reflectType.Field(i)
		currentValue := reflectValue.FieldByName(currentField.Name)
		multipart := currentField.Tag.Get("multipart")
		jsonField := currentField.Tag.Get("json")

		required := currentField.Tag.Get("required")
		requiredArray := strings.Split(required, ",")
		empty := currentField.Tag.Get("empty")

		switch multipart {
		case "file":
			if !m.basicValidator.ValidateStringContainInStringArray(requiredArray, param.menu) {
				continue
			}

			file, err := m.reader.ReadMultipartFile(
				param.r,
				jsonField,
			)

			if err != nil || len(file) == 0 {
				if empty != "allowed" {
					return nil, error2.ErrEmptyField.Param(jsonField).ContextModel(param.ctx)
				}

				continue
			}

			min, _ := strconv.Atoi(currentField.Tag.Get("min"))
			max, _ := strconv.Atoi(currentField.Tag.Get("max"))
			var totalData int64

			for i := 0; i < len(file); i++ {
				totalData += file[i].Header.Size

				if min > 0 && int(totalData) < min {
					return nil, error2.ErrFormatFieldRule.Param(jsonField, "NEED_MORE_THAN", strconv.Itoa(min)).ContextModel(param.ctx)
				}

				if max > 0 && int(totalData) > max {
					return nil, error2.ErrFormatFieldRule.Param(jsonField, "NEED_LESS_THAN", strconv.Itoa(max)).ContextModel(param.ctx)
				}

				bytes := param.ctx.ClientAccess.Logger.ModelLogger.ByteIn
				bytes += int(file[i].Header.Size)
				param.ctx.ClientAccess.Logger.Set(constanta.LoggerByteIn, bytes)

				extensions := strings.Split(currentField.Tag.Get("ext"), ",")
				extensionfile := strings.Split(file[i].Header.Filename, ".")

				if !m.basicValidator.ValidateStringContainInStringArray(
					extensions,
					extensionfile[len(extensionfile)-1],
				) {
					return nil, error2.ErrUnknownData.Param(jsonField).ContextModel(param.ctx)
				}

				if multipart == "file" {
					currentValue.Set(reflect.ValueOf(reflect.ValueOf(file[i]).Interface()))
					continue
				} else {
					currentValue.Set(reflect.ValueOf(reflect.ValueOf(file).Interface()))
				}
			}
		case "json":
			value := param.r.FormValue(jsonField)
			if !m.basicValidator.ValidateStringContainInStringArray(requiredArray, param.menu) {
				continue
			}

			if value == "" && empty == "allowed" {
				continue
			}

			bytes := param.ctx.ClientAccess.Logger.ModelLogger.ByteIn
			bytes += len(value)
			param.ctx.ClientAccess.Logger.Set(constanta.LoggerByteIn, bytes)

			if value == "" {
				return nil, error2.ErrEmptyField.Param(jsonField).ContextModel(param.ctx)
			}

			temp := reflect.New(currentField.Type)
			ptr := temp.Interface()
			err := json.Unmarshal([]byte(value), ptr)

			if err != nil {
				return nil, error2.ErrReadBody
			}

			currentValue.Set(reflect.ValueOf(ptr).Elem())
		default:
			value := param.r.FormValue(jsonField)
			bytes := param.ctx.ClientAccess.Logger.ModelLogger.ByteIn
			bytes += len(value)
			param.ctx.ClientAccess.Logger.Set(constanta.LoggerByteIn, bytes)

			switch currentValue.Kind() {
			case reflect.String:
				currentValue.SetString(value)
			case reflect.Int:
				if val, err := strconv.Atoi(value); err == nil {
					currentValue.SetInt(int64(val))
				} else {
					return nil, error2.ErrFormatField.Param(jsonField).ContextModel(param.ctx)
				}
			case reflect.Float32:
				if val, err := strconv.ParseFloat(value, 32); err == nil {
					currentValue.SetFloat(val)
				} else {
					return nil, error2.ErrFormatField.Param(jsonField).ContextModel(param.ctx)
				}
			case reflect.Float64:
				if val, err := strconv.ParseFloat(value, 64); err == nil {
					currentValue.SetFloat(val)
				} else {
					return nil, error2.ErrFormatField.Param(jsonField).ContextModel(param.ctx)
				}
			}
		}
	}

	err := m.tagValidator.ValidateByTag(param.ctx, param.dto, param.functionValidator, param.menu)

	if err != nil {
		return nil, err
	}

	return param.dto, nil
}
