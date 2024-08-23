package error

import "errors"

var ErrUnauthorized = NewUnBundledErrorMessages(401, errors.New("E-1-CMD-AUT-001"), nil)
var ErrExpiredToken = NewUnBundledErrorMessages(401, errors.New("E-1-CMD-AUT-002"), nil)

var ErrFieldInvalid = NewUnBundledErrorMessages(400, errors.New("E-4-GL-SRV-001"), errFieldNameConverter)
var ErrEmptyField = NewUnBundledErrorMessages(400, errors.New("E-4-CMD-DTO-001"), errFieldNameConverter)
var ErrFormatField = NewUnBundledErrorMessages(400, errors.New("E-4-CMD-DTO-002"), errFieldNameConverter)
var ErrFormatFieldRule = NewUnBundledErrorMessages(400, errors.New("E-4-CMD-DTO-003"), errFieldRuleConverter)
var ErrUnknownData = NewUnBundledErrorMessages(400, errors.New("E-4-CMD-DTO-004"), errFieldNameConverter)
var ErrDataUsed = NewUnBundledErrorMessages(400, errors.New("E-4-GL-SRV-004"), errFieldNameConverter)
var ErrReservedValueString = NewUnBundledErrorMessages(400, errors.New("E-4-CMD-DTO-006"), errFieldNameConverter)

var ErrReadBody = NewUnBundledErrorMessages(400, errors.New("E-4-CMD-BDY-001"), nil)
var ErrMarshalingBody = NewUnBundledErrorMessages(400, errors.New("E-4-CMD-BDY-002"), nil)

var ErrJWTForbiddenByResource = NewUnBundledErrorMessages(403, errors.New("E-3-CMD-AUT-001"), nil)
var ErrInvalidIPAddress = NewUnBundledErrorMessages(403, errors.New("E-3-CMD-AUT-004"), nil)
var ErrInvalidURLWhitelist = NewUnBundledErrorMessages(403, errors.New("E-3-CMD-AUT-005"), nil)

var errFieldNameConverter = func(value ...interface{}) map[string]ErrorParam {
	result := make(map[string]ErrorParam)
	result["FieldName"] = ErrorParam{value[0], true}
	return result
}

var errFieldRuleConverter = func(value ...interface{}) map[string]ErrorParam {
	result := make(map[string]ErrorParam)
	result["FieldName"] = ErrorParam{value[0], true}
	result["RuleName"] = ErrorParam{value[1], true}
	result["Other"] = ErrorParam{value[2], false}
	return result
}
