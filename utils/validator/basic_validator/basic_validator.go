package basic_validator

import "strings"

func (bv BasicValidator) CheckIsResourceIDExist(availableResourceID string, resourceIDMustHave string) bool {
	splitDot := strings.Split(availableResourceID, " ")
	return bv.ValidateStringContainInStringArray(splitDot, resourceIDMustHave)
}

func (bv BasicValidator) ValidateStringContainInStringArray(listString []string, key string) bool {
	for i := 0; i < len(listString); i++ {
		if listString[i] == key {
			return true
		}
	}
	return false
}
