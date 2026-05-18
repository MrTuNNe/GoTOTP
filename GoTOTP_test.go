package GoTOTP

import (
	"testing"
)

var secretKey = "OK6ZZOALZY6RNZBPM4QKD2ZFO5F3PTP56VIAXLDJLEHBPLJJIZNQ"

var totp = TOTP{Key: secretKey}

func TestTOTP_Verify(t *testing.T) {
	// this should fail as is it and old code expired (with the current timestamp)
	if totp.Verify("149425") {
		t.Error("Expected behavior is to fail. Check the implementation")
	}
}

func TestTOTP_VerifyTimestampOK(t *testing.T) {
	if !totp.VerifyWithTimestamp(1723719527, "611626") { // this should verify as good
		t.Error("Expected behavior is to accept the code. Check the implementation")
	}
	if totp.VerifyWithTimestamp(1723719580, "611626") { // past the 30 seconds, supposed to fail
		t.Error("Expected behavior is to fail. Check the implementation")
	}
}

func TestTOTP_GenerateTOTP(t *testing.T) {
	code, err := totp.GenerateTOTP(1723719527) // should generate a code (611626), not to fail
	if err != nil {
		t.Error("Generating the code has failed. Was supposed to work.")
	}
	if code != "611626" {
		t.Error("The generated code is supposed to be `611626` but it's not")
	}
}

func TestTOTP_RandomSecret(t *testing.T) {
	secret, err := GenerateRandomSecret(32)
	if err != nil {
		t.Error(err)
	}
	totp_test := TOTP{Key: secret}
	_, err = totp_test.GenerateTOTP(1723719527)
	if err != nil {
		t.Error("Generating the code has failed. The secret key might be problematic")
	}
}

func TestTOTP_SHA1(t *testing.T) {
	totp_sha1 := TOTP{Key: secretKey, Algorithm: SHA1}
	code, err := totp_sha1.GenerateTOTP(1723719527)
	if err != nil {
		t.Error("Generating the code has failed. Was supposed to work.")
	}
	if code != "032332" {
		t.Errorf("Expected `032332` but got `%s`", code)
	}
	if !totp_sha1.VerifyWithTimestamp(1723719527, "032332") {
		t.Error("Expected SHA1 code to be accepted")
	}
	if totp_sha1.VerifyWithTimestamp(1723719580, "032332") {
		t.Error("Expected SHA1 code to be rejected past its window")
	}
	if totp_sha1.VerifyWithTimestamp(1723719527, "611626") {
		t.Error("Expected SHA256 code to be rejected when algorithm is SHA1")
	}
}

func TestTOTP_SHA512(t *testing.T) {
	totp_sha512 := TOTP{Key: secretKey, Algorithm: SHA512}
	code, err := totp_sha512.GenerateTOTP(1723719527)
	if err != nil {
		t.Error("Generating the code has failed. Was supposed to work.")
	}
	if code != "319711" {
		t.Errorf("Expected `319711` but got `%s`", code)
	}
	if !totp_sha512.VerifyWithTimestamp(1723719527, "319711") {
		t.Error("Expected SHA512 code to be accepted")
	}
	if totp_sha512.VerifyWithTimestamp(1723719580, "319711") {
		t.Error("Expected SHA512 code to be rejected past its window")
	}
	if totp_sha512.VerifyWithTimestamp(1723719527, "611626") {
		t.Error("Expected SHA256 code to be rejected when algorithm is SHA512")
	}
}

func TestTOTP_AlgorithmIsolation(t *testing.T) {
	// All three algorithms must produce distinct codes for the same key/timestamp
	ts := int64(1723719527)
	codes := map[Algorithm]string{}
	for _, alg := range []Algorithm{SHA1, SHA256, SHA512} {
		code, err := (&TOTP{Key: secretKey, Algorithm: alg}).GenerateTOTP(ts)
		if err != nil {
			t.Errorf("%s: unexpected error: %v", alg, err)
		}
		codes[alg] = code
	}
	if codes[SHA1] == codes[SHA256] || codes[SHA256] == codes[SHA512] || codes[SHA1] == codes[SHA512] {
		t.Errorf("Expected distinct codes per algorithm, got: %v", codes)
	}
}

func TestTOTP_GenerateURI(t *testing.T) {
	otp_good := TOTP{
		Key:      secretKey,
		Issuer:   "mrtunne.info",
		UserName: "admin@admin.test",
	}
	_, err := otp_good.GenerateURI()
	if err != nil {
		t.Error(err)
	}
	otp_bad := TOTP{
		Key: secretKey,
	}
	_, err = otp_bad.GenerateURI()
	if err == nil {
		t.Error("This implementation was supposed to return errors as it has null values for `Issuer` and `UserName`")
	}
}
