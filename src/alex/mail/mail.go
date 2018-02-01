package mail

import (
    "gopkg.in/gomail.v2"
)

type Mail struct {
    Host     string
    Port     int
    From     string
    Password string
}

func New(host string, port int, from string, password string) *Mail {
    m := new(Mail)
    m.Host = host
    m.Port = port
    m.From = from
    m.Password = password
    return m
}

//发送邮件
//@param title 标题
//@param html 内容
//@param to 发送对象
func (t *Mail) Send(title string, html string, to []string) error {
    m := gomail.NewMessage()
    m.SetHeader("From", t.From)
    m.SetHeader("To", to...)
    m.SetHeader("Subject", title)
    m.SetBody("text/html", html)
    d := gomail.NewDialer(t.Host, t.Port, t.From, t.Password)
    if err := d.DialAndSend(m); err != nil {
        return err
    }
    return nil
}
