package GoTOTP

import (
	"crypto/hmac"
	"crypto/sha1"
	"crypto/sha256"
	"crypto/sha512"
	"encoding/base32"
	"encoding/binary"
	"errors"
	"fmt"
	"hash"
	"math"
	"net/url"
	"strconv"
	"strings"
	"time"
)

type Algorithm string

const (
	SHA1   Algorithm = "SHA1"
	SHA256 Algorithm = "SHA256"
	SHA512 Algorithm = "SHA512"
)

type TOTP struct {
	Key       string
	Issuer    string
	UserName  string
	Algorithm Algorithm
}

func (totp *TOTP) hmacHash(message []byte) ([]byte, error) {
	key, err := totp.validateSecret()
	if err != nil {
		return []byte{}, err
	}
	var h func() hash.Hash
	switch totp.Algorithm {
	case SHA1:
		h = sha1.New
	case SHA512:
		h = sha512.New
	default:
		h = sha256.New
	}
	mac := hmac.New(h, key)
	mac.Write(message)
	return mac.Sum(nil), nil
}

func (totp *TOTP) validateSecret() ([]byte, error) {
	key := totp.Key
	if len(key)%8 != 0 {
		key = key + strings.Repeat("=", 8-(len(key)%8))
	}
	return base32.StdEncoding.DecodeString(key)
}

// Based from RFC 6238
func (totp *TOTP) GenerateTOTP(timestamp int64) (string, error) {
	if totp.Key == "" {
		return "", errors.New("`Key` value cannot be empty")
	}
	codeDigits := 6
	var result string
	currentTime := timestamp / int64(30)
	buf := make([]byte, 8)
	binary.BigEndian.PutUint64(buf, uint64(currentTime))
	hash, err := totp.hmacHash(buf)
	if err != nil {
		return "", err
	}
	offset := int(hash[len(hash)-1] & 0xf)
	code := (int(hash[offset]&0x7f) << 24) |
		(int(hash[offset+1]&0xff) << 16) |
		(int(hash[offset+2]&0xff) << 8) |
		(int(hash[offset+3] & 0xff))
	code = code % int(math.Pow10(codeDigits))
	result = strconv.Itoa(code)
	for len(result) < codeDigits {
		result = "0" + result
	}
	return result, nil
}

// Verify if the given input code is valid for the current timestamp.
func (totp *TOTP) Verify(inputCode string) bool {
	timestamp := time.Now().Unix()
	code, err := totp.GenerateTOTP(timestamp)
	if err != nil {
		return false
	}
	return code == inputCode
}

// Verify if the input code is valid for a given timestamp.
func (totp *TOTP) VerifyWithTimestamp(timestamp int64, inputCode string) bool {
	code, err := totp.GenerateTOTP(timestamp)
	if err != nil {
		return false
	}
	return code == inputCode
}

func (totp *TOTP) GenerateURI() (string, error) {
	if totp.Issuer == "" || totp.UserName == "" || totp.Key == "" {
		return "", errors.New("you must specify a value for `Issuer`, `UserName` and `Key` to generate an URI")
	}
	algorithm := totp.Algorithm
	if algorithm == "" {
		algorithm = SHA256
	}
	uri := url.URL{
		Scheme: "otpauth",
		Host:   "totp",
		Path:   fmt.Sprintf("%s: %s", totp.Issuer, totp.UserName),
	}
	q := uri.Query()
	q.Add("secret", totp.Key)
	q.Add("issuer", totp.Issuer)
	q.Add("algorithm", string(algorithm))
	q.Add("digits", "6")
	q.Add("period", "30")
	uri.RawQuery = q.Encode()
	return uri.String(), nil
}
