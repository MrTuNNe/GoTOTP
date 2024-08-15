package totp

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base32"
	"encoding/binary"
	"log"
	"math"
	"strconv"
	"strings"
	"time"
)

type TOTP struct {
	Key string
}

func (totp *TOTP) hmac_sha256(message []byte) ([]byte, error) {
	key, err := totp.validateSecret()
	if err != nil {
		return []byte{}, err
	}
	mac := hmac.New(sha256.New, key)
	mac.Write(message)
	return mac.Sum(nil), nil
}

func (totp *TOTP) validateSecret() ([]byte, error) {
	// we add padding to the base32 secret key if necessary
	if len(totp.Key)%8 != 0 {
		totp.Key = totp.Key + strings.Repeat("=", 8-(len(totp.Key)%8))
	}
	return base32.StdEncoding.DecodeString(totp.Key)
}

// Based from RFC 6238
func (totp *TOTP) GenerateTOTP(timestamp int64) string {
	codeDigits := 6
	var result string
	currentTime := timestamp / int64(30)
	// we convert the timestamp from int64 to a byte array
	buf := make([]byte, 8)
	binary.BigEndian.PutUint64(buf, uint64(currentTime))
	hash, err := totp.hmac_sha256(buf)
	if err != nil {
		log.Fatal(err)
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
	return result
}

// Verify if the given input code is valid for the current timestamp.
func (totp *TOTP) Verify(inputCode string) bool {
	timestamp := time.Now().Unix()
	code := totp.GenerateTOTP(timestamp)
	return code == inputCode
}

// Verify if the input code is valid for a given timestamp. Use this just for testing
func (totp *TOTP) VerifyWithTimestamp(timestamp int64, inputCode string) bool {
	code := totp.GenerateTOTP(timestamp)
	return code == inputCode
}
