package router

import (
	"crypto/hmac"
	"crypto/sha256"
	"crypto/sha512"
	"encoding/hex"
	"hash"
	"strings"
)

func GenerateMessageDigest(
	message string,
) string {
	message = strings.Replace(message, "\r", "", -1)
	return CheckSumWithSha256([]byte(message))
}

func CheckSumWithSha256(content []byte) string {
	result := sha256.Sum256(content)
	return hex.EncodeToString(result[:])
}

func GenerateSignature(
	httpMethod string,
	relativeURL string,
	accessToken string,
	messageDigest string,
	timestamp string,
	key string,
) string {
	message := httpMethod + ":" + relativeURL + ":" + accessToken + ":" + messageDigest + ":" + timestamp
	return ChecksumWithHMACSHA(sha512.New, []byte(message), key)
}

func ChecksumWithHMACSHA(f func() hash.Hash, content []byte, key string) string {
	h := hmac.New(f, []byte(key))
	h.Write(content)
	return hex.EncodeToString(h.Sum(nil))
}
