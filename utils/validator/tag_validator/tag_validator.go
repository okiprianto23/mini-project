package tag_validator

import (
	"main-xyz/constanta"
	"main-xyz/context"
	error2 "main-xyz/error"
	rgx "main-xyz/regex"
	"main-xyz/utils/validator/basic_validator"
	"reflect"
	"regexp"
	"strconv"
	"strings"
	"time"
)

func NewTagValidator() TagValidator {
	validator := &tagValidator{
		basic: basic_validator.BasicValidator{},
	}
	validator.build()
	return validator
}

type tagValidator struct {
	enum       map[string][]string
	dateFormat map[string]string
	regex      map[string]regexValue
	basic      basic_validator.BasicValidator
}

func (v *tagValidator) AddEnum(key string, value ...string) *tagValidator {
	v.enum[key] = value
	return v
}

func (v *tagValidator) AddDateFormat(key string, dateFormat string) *tagValidator {
	v.dateFormat[key] = dateFormat
	return v
}

func (v *tagValidator) AddRegex(key string, regex string, name string) *tagValidator {
	v.regex[key] = regexValue{
		regex:    regex,
		ruleName: name,
	}
	return v
}

func (v *tagValidator) build() {
	v.enum = make(map[string][]string)
	v.enum["boolean_permission"] = []string{"Y", "N"}

	v.dateFormat = make(map[string]string)
	v.dateFormat["default"] = constanta.DefaultTimeFormat
	v.dateFormat["date_only"] = constanta.DateOnlyTimeFormat

	v.regex = make(map[string]regexValue)
	v.regex["profile_name"] = regexValue{
		regex:    rgx.PROFILE_NAME,
		ruleName: "PROFILE_NAME_REGEX_MESSAGE",
	}
	v.regex["directory_name"] = regexValue{
		regex:    rgx.DIRECTORY_NAME,
		ruleName: "DIRECTORY_NAME_REGEX_MESSAGE",
	}
	v.regex["name"] = regexValue{
		regex:    rgx.NAME_STANDARD,
		ruleName: "NAME_REGEX_MESSAGE",
	}
	v.regex["text_only"] = regexValue{
		regex:    rgx.TEXT_ONLY,
		ruleName: "DESCRIPTION_REGEX_MESSAGE",
	}
	v.regex["alphanumeric"] = regexValue{
		regex:    rgx.ALPHANUMERIC,
		ruleName: "ALPHANUMERIC_REGEX",
	}
	v.regex["country_code"] = regexValue{
		regex:    rgx.COUNTRY_CODE,
		ruleName: "COUNTRY_CODE_REGEX",
	}
	v.regex["email"] = regexValue{
		regex:    rgx.EMAIL_REGEX,
		ruleName: "EMAIL_REGEX",
	}
	v.regex["phone"] = regexValue{
		regex:    rgx.PHONE_NUMBER_WITH_COUNTRY_CODE,
		ruleName: "PHONE_NUMBER_REGEX",
	}
	v.regex["username"] = regexValue{
		regex:    rgx.USERNAME,
		ruleName: "USERNAME_REGEX",
	}
	v.regex["lowercase"] = regexValue{
		regex:    rgx.LOWERCASE,
		ruleName: "LOWERCASE_REGEX",
	}
	v.regex["upercase"] = regexValue{
		regex:    rgx.UPERCASE,
		ruleName: "UPERCASE_REGEX",
	}
	v.regex["nik"] = regexValue{
		regex:    rgx.NIK,
		ruleName: "NIK_REGEX",
	}
	v.regex["lowercase_number"] = regexValue{
		regex:    rgx.LOWERCASE_AND_NUMBER,
		ruleName: "LOWERCASE_AND_NUMBER_REGEX",
	}
	v.regex["numeric"] = regexValue{
		regex:    rgx.LONG_NUMERIC,
		ruleName: "LONG_NUMERIC_REGEX",
	}
}

var (
	intType   = []string{reflect.Int64.String(), reflect.Int.String(), reflect.Int32.String()}
	floatType = []string{reflect.Float32.String(), reflect.Float64.String()}
)

func (v *tagValidator) ValidateByTag(
	ctx *context.ContextModel,
	dto interface{},
	functionValidator string,
	menu string,
) (
	err error,
) {

	reflectType := reflect.TypeOf(dto).Elem()
	reflectValue := reflect.ValueOf(dto).Elem()

	max := 0
	min := 0
	isMinFound := false
	isMaxFound := false

	for i := 0; i < reflectType.NumField(); i++ {
		currentField := reflectType.Field(i)
		currentValue := reflectValue.FieldByName(currentField.Name)

		// ini untuk case jika menggunakan form data agar bisa check validasi yg strucknya
		if currentField.Type.Kind() == reflect.Struct {
			newDTO := currentValue.Addr().Interface()
			err = v.ValidateByTag(ctx, newDTO, functionValidator, menu)
			if err != nil {
				return
			}
		}

		//Untuk merubah nilai agar sesuai dengan yg di minta
		//Contoh uppercase nanti valuenya akan secara otomatis jadi upplercase
		if currentField.Type.Kind() == reflect.String {
			autoFix := currentField.Tag.Get("auto_fix")
			split := strings.Split(autoFix, ",")
			for i := 0; i < len(split); i++ {
				if split[i] == "uppercase" {
					currentValue.SetString(strings.ToUpper(currentValue.String()))
				} else if split[i] == "lowercase" {
					currentValue.SetString(strings.ToLower(currentValue.String()))
				} else if split[i] == "trim" {
					currentValue.SetString(strings.Trim(currentValue.String(), " "))
				} else if split[i] == "filename" {
					val := currentValue.String()
					val = strings.ReplaceAll(val, "\\", "/")
					val = strings.Trim(val, "./")
					val = strings.Trim(val, "/")
					currentValue.SetString(val)
				}
			}
		}

		//for validate required
		required := currentField.Tag.Get("required")
		jsonField := strings.ToUpper(currentField.Tag.Get("json"))
		requiredArray := strings.Split(required, ",")

		// boleh kosong
		empty := currentField.Tag.Get("empty")
		if empty == "allowed" {
			if currentValue.IsZero() && currentField.Type.String() != "time.Time" {
				continue
			}
		}

		minEqMax := false

		if v.basic.ValidateStringContainInStringArray(requiredArray, menu) {
			defaultValue := currentField.Tag.Get("default")
			min, isMinFound, max, isMaxFound = v.getMinMaxValue(currentField)
			minEqMax = min == max

			if v.basic.ValidateStringContainInStringArray(intType, currentField.Type.String()) {
				if currentValue.IsZero() {
					valueIn, _ := strconv.Atoi(defaultValue)
					currentValue.SetInt(int64(valueIn))
				}

				value := currentValue.Int()
				if isMinFound {
					if min != 0 && int(value) == 0 {
						err = error2.ErrEmptyField.Param(jsonField).ContextModel(ctx)
						return
					}
					if int(value) < min {
						return error2.ErrFormatFieldRule.Param(jsonField, getRuleName(minEqMax, "NEED_MORE_THAN"), strconv.Itoa(min)).ContextModel(ctx)
					}
				}
				if isMaxFound {
					if int(value) > max {
						return error2.ErrFormatFieldRule.Param(jsonField, getRuleName(minEqMax, "NEED_LESS_THAN"), strconv.Itoa(max)).ContextModel(ctx)
					}
				}
			} else if v.basic.ValidateStringContainInStringArray(floatType, currentField.Type.String()) {
				if currentValue.IsZero() {
					valueIn, _ := strconv.ParseFloat(defaultValue, 64)
					currentValue.SetFloat(valueIn)
				}

				value := currentValue.Float()
				if isMinFound {
					if value < float64(min) {
						return error2.ErrFormatFieldRule.Param(jsonField, getRuleName(minEqMax, "NEED_MORE_THAN"), strconv.Itoa(min)).ContextModel(ctx)
					}
				}
				if isMaxFound {
					if value > float64(max) {
						return error2.ErrFormatFieldRule.Param(jsonField, getRuleName(minEqMax, "NEED_LESS_THAN"), strconv.Itoa(max)).ContextModel(ctx)
					}
				}
			} else if reflect.String.String() == currentField.Type.String() {
				currentValue.SetString(strings.Trim(currentValue.String(), " "))
				if currentValue.IsZero() {
					currentValue.SetString(defaultValue)
				}

				if currentValue.IsZero() {
					return error2.ErrEmptyField.Param(jsonField).ContextModel(ctx)
				}

				value := currentValue.String()
				err = v.ValidateMinMaxString(ctx, value, jsonField, min, max)
				if err != nil {
					return
				}

				enumField := currentField.Tag.Get("enum")
				if enumField != "" {
					if !v.basic.ValidateStringContainInStringArray(v.enum[enumField], currentValue.String()) {
						err = error2.ErrFormatFieldRule.Param(jsonField, "FIXED_VALUE", strings.Join(v.enum[enumField], " , ")).ContextModel(ctx)
						return
					}
				}

				regexField := currentField.Tag.Get("regex")
				if regexField != "" {
					if v.regex[regexField].regex != "" {
						if len(value) > 0 || min != 0 {
							if !regexp.MustCompile(v.regex[regexField].regex).MatchString(currentValue.String()) {
								return error2.ErrFormatFieldRule.Param(jsonField, v.regex[regexField].ruleName, "").ContextModel(ctx)
							}
						}

					}
				}
			} else if "time.Time" == currentField.Type.String() {
				var timeObject time.Time
				dateFormatTag := currentField.Tag.Get("dateFormat")
				strField := currentField.Name + "Str"
				timeFormatUsed := v.dateFormat["default"]

				dateSplit := strings.Split(dateFormatTag, ",")
				var timeValid = false

				for i := 0; i < len(dateSplit); i++ {
					if v.dateFormat[dateSplit[i]] != "" {
						timeFormatUsed = v.dateFormat[dateSplit[i]]
					}

					field, valid := reflectType.FieldByName(strField)

					if valid {

						val := reflectValue.FieldByName(strField).String()
						if val == "" && empty == "allowed" {
							timeValid = true
							continue
						}

						jsonField = field.Tag.Get("json")
						timeObject, err = v.TimeStrToTime(ctx, val, jsonField, timeFormatUsed)

						if err != nil {
							continue
						}

						currentValue.Set(reflect.ValueOf(timeObject))
						timeValid = true
						break

					} else {
						return error2.ErrUnknownData.Param(strField).ContextModel(ctx)
					}
				}

				if !timeValid {
					return error2.ErrFormatField.Param(jsonField).ContextModel(ctx)
				}
			} else if currentValue.Kind() == reflect.Slice {
				if isMinFound {
					if currentValue.Len() == 0 {
						return error2.ErrEmptyField.Param(jsonField).ContextModel(ctx)
					}
					if currentValue.Len() < min {
						return error2.ErrFormatFieldRule.Param(jsonField, getRuleName(minEqMax, "NEED_MORE_THAN"), strconv.Itoa(min)).ContextModel(ctx)
					}
				}

				if isMaxFound {
					if currentValue.Len() > max {
						return error2.ErrFormatFieldRule.Param(jsonField, getRuleName(minEqMax, "NEED_LESS_THAN"), strconv.Itoa(max)).ContextModel(ctx)
					}
				}

				for i := 0; i < currentValue.Len(); i++ {
					temp := currentValue.Index(i)
					if temp.Type().String() == reflect.Struct.String() || temp.Type().Kind().String() == reflect.Struct.String() {
						newDTO := currentValue.Index(i).Addr().Interface()
						err = v.ValidateByTag(ctx, newDTO, functionValidator, menu)
						if err != nil {
							return
						}
					}

					if temp.Type().Kind() == reflect.String {
						enumField := currentField.Tag.Get("enum")
						if enumField != "" {
							if !v.basic.ValidateStringContainInStringArray(v.enum[enumField], temp.String()) {
								err = error2.ErrFormatFieldRule.Param(jsonField, "FIXED_VALUE", strings.Join(v.enum[enumField], " , ")).ContextModel(ctx)
								return
							}
						}
					}
				}
			}
		}
	}

	return
}

func getRuleName(minEqmax bool, ruleName string) string {
	if minEqmax {
		return "EQUAL"
	}

	return ruleName
}

func (v *tagValidator) getMinMaxValue(field reflect.StructField) (min int, isMinFound bool, max int, isMaxFound bool) {
	maxStr, isMaxFound := field.Tag.Lookup("max")
	minStr, isMinFound := field.Tag.Lookup("min")

	min, _ = strconv.Atoi(minStr)
	max, _ = strconv.Atoi(maxStr)

	return
}

func (v *tagValidator) TimeStrToTime(
	ctx *context.ContextModel,
	timeStr string,
	fieldName string,
	format string,
) (
	output time.Time,
	err error,
) {
	output, err = time.Parse(format, timeStr)
	if err != nil {
		err = error2.ErrFormatField.Param(fieldName).ContextModel(ctx)
		return
	}

	return output, nil
}

func (v *tagValidator) ValidateMinMaxString(
	ctx *context.ContextModel,
	inputStr string,
	fieldName string,
	min int,
	max int,
) error {
	minEqMax := min == max
	if min != 0 {
		if len(inputStr) == 0 {
			return error2.ErrEmptyField.Param(fieldName).ContextModel(ctx)
		}
		if len(inputStr) < min {
			if min == 1 {
				return error2.ErrEmptyField.Param(fieldName).ContextModel(ctx)
			} else {
				return error2.ErrFormatFieldRule.Param(fieldName, getRuleName(minEqMax, "NEED_MORE_THAN"), strconv.Itoa(min)).ContextModel(ctx)
			}
		}
	}
	if max != 0 {
		if len(inputStr) > max {
			return error2.ErrFormatFieldRule.Param(fieldName, getRuleName(minEqMax, "NEED_LESS_THAN"), strconv.Itoa(max)).ContextModel(ctx)
		}
	}

	return nil
}
