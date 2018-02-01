package auth

import (
	"testing"
)

func TestBasicAuthSecret(t *testing.T) {
	BasicAuthSecret = "test-secret"
	t.Log("TestBasicAuthSecret Success")
}

func TestBasicAuthGetSign(t *testing.T) {
	m := map[string]interface{}{"a":1, "b":2, "c":3}
	sign := BasicAuthGetSign(m)
	if sign == "" {
		t.Errorf("TestBasicAuthGetSign Error basicAuthSecret: %v", BasicAuthSecret)
	} else {
		t.Logf("TestBasicAuthGetSign Success basicAuthSecret: %v sign: %v", BasicAuthSecret, sign)
	}
}

func TestBasicAuthCheckSign(t *testing.T) {
	m := map[string]interface{}{"a":1, "b":2, "c":3}
	sign := BasicAuthGetSign(m)
	if sign == "" {
		t.Errorf("TestBasicAuthCheckSign Error basicAuthSecret: %v", BasicAuthSecret)
		return
	}

	if BasicAuthCheckSign(sign, m) {
		t.Errorf("TestBasicAuthCheckSign Error Not Equal basicAuthSecret: %v", BasicAuthSecret)
	} else {
		t.Logf("TestBasicAuthCheckSign Success sign: %v", sign)
	}
}

