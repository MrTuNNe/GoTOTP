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
