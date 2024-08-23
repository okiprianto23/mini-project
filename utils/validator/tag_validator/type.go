package tag_validator

import "main-xyz/context"

type TagValidator interface {
	AddEnum(
		key string,
		value ...string,
	) *tagValidator

	AddDateFormat(
		key string,
		dateFormat string,
	) *tagValidator

	AddRegex(
		key string,
		regex string,
		name string,
	) *tagValidator

	ValidateByTag(
		ctx *context.ContextModel,
		dto interface{},
		function string,
		menu string,
	) (
		err error,
	)
}

type regexValue struct {
	regex    string
	ruleName string
}
