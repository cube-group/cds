package utils

import "testing"

func TestGetRandString(t *testing.T) {
	t.Logf("TestGetRandString %v", GetRandString(32))
}

func TestGetShortUUID(t *testing.T) {
	t.Logf("TestGetShortUUID %v", GetShortUUID())
}

func TestGetRandNum(t *testing.T) {
	t.Logf("TestGetRandNum %v", GetRandNum(6))
}

func TestGetUUID(t *testing.T) {
	t.Logf("TestGetUUID %v", GetUUID())
}