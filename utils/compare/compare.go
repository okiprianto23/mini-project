package compare

// IsEmptyString checks if a given string is empty.
func IsEmptyString(s string) bool {
	return s == ""
}

// IsEmptyInt64 checks if a given int64 is empty.
func IsEmptyInt64(i int64) bool {
	return i < 1
}
