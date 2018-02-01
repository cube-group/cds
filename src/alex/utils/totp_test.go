package utils

import (
	"testing"
	"fmt"
)

const (
	TOTP_USERNAME = "fcds"
	TOTP_MAIL = "lin2798003@sina.com"
)

func TestTotpCode(t *testing.T) {
	equal, err := TotpCode("566781", "RYBTGKIAZROTYRTT")
	fmt.Println(equal, err)
}

func TestTotpUrl(t *testing.T) {
	url, secret, err := TotpUrlAndSecret(TOTP_USERNAME, TOTP_MAIL)
	fmt.Println(url, secret, err)
}
