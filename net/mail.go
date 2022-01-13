package net

import (
	"errors"
	"fmt"
	"net/smtp"
	"strings"

	"github.com/xurwxj/gtils/base"
	"github.com/xurwxj/viper"
)

// SendMail need following config in config.json:
// "email": {
// 		"account": "xxx@xxx.com",
// 		"pwd": "xxx",
// 		"host": "smtpdm.xxx.com:25"
//   },
func SendMail(to, fromUser, subject, body, mailType string, options ...string) error {
	if to != "" && subject != "" {
		user := viper.GetString("email.account")
		password := viper.GetString("email.pwd")
		host := viper.GetString("email.host")
		for i, o := range options {
			if !base.FindInStringSlice([]string{"", "default"}, o) && i == 0 {
				// fmt.Println(o)
				user = viper.GetString(fmt.Sprintf("email.%s.account", o))
				password = viper.GetString(fmt.Sprintf("email.%s.pwd", o))
				host = viper.GetString(fmt.Sprintf("email.%s.host", o))
			}
		}
		// fmt.Println(user)
		if password == "" || host == "" || user == "" {
			return fmt.Errorf("authParamsErr")
		}
		auth := LoginAuth(user, password)
		var contentType string
		if mailType == "html" {
			contentType = "Content-Type: text/" + mailType + "; charset=UTF-8"
		} else {
			contentType = "Content-Type: text/plain" + "; charset=UTF-8"
		}
		if fromUser == "" {
			fromUser = user
		}
		msg := []byte("To: " + to + "\r\nFrom: " + fromUser + "\r\nSubject: " + subject + "\r\n" + contentType + "\r\n\r\n" + body)
		sendTo := strings.Split(to, ";")
		if len(sendTo) > 0 && subject != "" {
			err := smtp.SendMail(host, auth, user, sendTo, msg)
			return err
		}
		return nil
	}
	return errors.New("paramsErr")
}

type loginAuth struct {
	username, password string
}

// LoginAuth login auth for mail send
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
