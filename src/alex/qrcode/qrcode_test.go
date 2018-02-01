package qrcode

import "testing"

func TestPng(t *testing.T) {

    err := Png(
        "otpauth://totp/fcds:cqh?algorithm=SHA1&digits=6&issuer=fcds&period=30&secret=K25HNAAJR7WVSVSN",
        "qrcodeTest.png",
        200,
        200,
    )
    if err != nil {
        t.Error(err)
    }

}