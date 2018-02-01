package utils

import (
	"github.com/pquerna/otp/totp"
)

//校验totp正确性
func TotpCode(code, secret string) (bool, error) {
	return totp.Validate(code, secret), nil
}

//返回totp url和secret
func TotpUrlAndSecret(issuer, mail string) (string, string, error) {
	key, err := totp.Generate(totp.GenerateOpts{Issuer:issuer, AccountName:mail})
	if err != nil {
		return "", "", err
	}
	return key.URL(), key.Secret(), nil
}