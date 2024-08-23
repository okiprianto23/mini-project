package utils

import (
	"main-xyz/constanta"
	error2 "main-xyz/error"
	"main-xyz/router"
	"strconv"
)

func CheckIDParam(param router.URLParam) (id int, err error) {
	if param.Path["ID"] == "null" {
		err = error2.ErrUnknownData.Param(constanta.ID)
		return
	}

	id, err = strconv.Atoi(param.Path["ID"])
	if err != nil {
		err = error2.ErrFieldInvalid.Param(constanta.ID)
		return
	}
	return
}
