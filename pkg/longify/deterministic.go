package longify

import (
	"crypto/sha256"
	"encoding/base64"
	"errors"
	"strings"
	"unsafe"
)

const hashLen = sha256.Size

var base64Encoding = base64.URLEncoding

func ForwardDeterminisic(s string) (string, error) {
	var b strings.Builder
	b.WriteString(s)
	var hash [hashLen]byte = sha256.Sum256([]byte(s))
	_, err := b.Write(hash[:])
	if err != nil {
		return "", err
	}
	str := b.String()
	result := base64Encoding.EncodeToString([]byte(str))
	return result, nil
}

var ErrTypo error = errors.New("has typo")
var ErrTooShort error = errors.New("too short")
var ErrStrangeHash error = errors.New("strange hash")

func BackwardDeterministic(s string) (string, error) {
	link, err := base64Encoding.DecodeString(s)
	if err != nil {
		return "", ErrTypo
	}
	if len(link) <= hashLen {
		return "", ErrTooShort
	}
	// hash := link[len(link)-hashLen:]
	// if !isValidHash(hash) {
	// 	return "", ErrStrangeHash
	// }
	ans := link[:len(link)-hashLen]
	return unsafe.String(&ans[0], len(ans)), nil
}

func isValidHashSum(hashSum []byte) bool {
	if len(hashSum) != hashLen {
		return false
	}
	for i := 0; i < len(hashSum); i += 1 {
		if '0' <= hashSum[i] && hashSum[i] <= '9' {
			continue
		}
		if 'a' <= hashSum[i] && hashSum[i] <= 'f' {
			continue
		}
		return false
	}

	return true
}
