package mail

import (
    "testing"
)

func TestMail_Send(t *testing.T) {
    m := New("smtp.exmail.qq.com", 465, "system@foryou56.com", "Foryou56**")
    err := m.Send("test send", "Hello <b>Light</b> weight <i>baby</i>! <hr><hr><hr>yeah buddy",
        []string{"chenqionghe@foryou56.com"})
    if err != nil {
        t.Error(err)
    } else {
        t.Log("send ok")
    }
}
