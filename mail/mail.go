package mail

import (
	"net/smtp"
	"errors"
	"fmt"
)

type Mail struct {
	Host    string
	Post    int
	Username string
	Password string
	to []string
}

func (m *Mail)AddAddr(to string){
	if m.to == nil {
		m.to = make([]string, 0)
	}
	m.to = append(m.to, to)
}

func (m *Mail) Send(title, body string) error {
	if len(title) == 0 && len(body) == 0 {
		return errors.New("Mail.Run title and body is nil")
	}
	var auth smtp.Auth
	if len(m.Username) > 0 {
		auth = smtp.PlainAuth("", m.Username, m.Password, m.Host)
	}
	post := 25
	if m.Post > 0 {
		post = m.Post
	}
	addr := fmt.Sprintf("%s:%v", m.Host, post)
	err := smtp.SendMail(addr, auth, title, m.to, []byte(body))
	if err != nil {
		return err
	}
	return nil
}
