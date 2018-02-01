package utils

import (
	"testing"
	"time"
)

func TestGetTimer(t *testing.T) {
	t.Logf("TestGetTimer: %v", GetTimer())
}

func TestGetMicroTimer(t *testing.T) {
	t.Logf("TestGetMicroTimer: %v", GetMicroTimer())
}

func TestGetFormatYmdHis(t *testing.T) {
	t.Logf("TestGetFormatYmdHis: %v", GetFormatYmdHis())
}

func TestGetFormatYmdHisByTime(t *testing.T) {
	t.Logf("TestGetFormatYmdHisByTime: %v", GetFormatYmdHisByTime(time.Now()))
}

func TestGetFormatYmdHisByUnix(t *testing.T) {
	t.Logf("TestGetFormatYmdHisByUnix: %v", GetFormatYmdHisByUnix(time.Now().Unix()))
}

func TestGetFormatYmd(t *testing.T) {
	t.Logf("TestGetFormatYmdByTime: %v", GetFormatYmd())
}

func TestGetFormatYmdByTime(t *testing.T) {
	t.Logf("TestGetFormatYmdByTime: %v", GetFormatYmdByTime(time.Now()))
}

func TestGetFormatYmdByUnix(t *testing.T) {
	t.Logf("TestGetFormatYmdByUnix: %v", GetFormatYmdByUnix(time.Now().Unix()))
}


