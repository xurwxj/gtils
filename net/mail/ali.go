package mail

import (
	"errors"
	"fmt"
	"net/smtp"
	"strings"

	"github.com/xurwxj/viper"
)

func AliMail(to, fromUser, subject, body, mailtype string) error {
	if to != "" && subject != "" {
		user := viper.GetString("email.account")
		password := viper.GetString("email.pwd")
		host := viper.GetString("email.host")
		if password == "" || host == "" || user == "" {
			return fmt.Errorf("config not ready!")
		}
		auth := LoginAuth(user, password)
		var content_type string
		if mailtype == "html" {
			content_type = "Content-Type: text/" + mailtype + "; charset=UTF-8"
		} else {
			content_type = "Content-Type: text/plain" + "; charset=UTF-8"
		}
		if fromUser == "" {
			fromUser = user
		}
		msg := []byte("To: " + to + "\r\nFrom: " + fromUser + "\r\nSubject: " + subject + "\r\n" + content_type + "\r\n\r\n" + body)
		sendTo := strings.Split(to, ";")
		if len(sendTo) > 0 && subject != "" {
			err := smtp.SendMail(host, auth, user, sendTo, msg)
			return err
		}
		return nil
	}
	return errors.New("no receiver")
}

type loginAuth struct {
	username, password string
}

func LoginAuth(username, password string) smtp.Auth {
	return &loginAuth{username, password}
}

func (a *loginAuth) Start(server *smtp.ServerInfo) (string, []byte, error) {
	return "LOGIN", []byte(a.username), nil
}

func (a *loginAuth) Next(fromServer []byte, more bool) ([]byte, error) {
	if more {
		switch string(fromServer) {
		case "Username:":
			return []byte(a.username), nil
		case "Password:":
			return []byte(a.password), nil
		}
	}
	return nil, nil
}
