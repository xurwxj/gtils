package mail

import "github.com/xurwxj/viper"

func SendMail(to, from, title, body, mailType string) error {
	cloud := viper.GetString("email.cloud")
	if cloud == "" {
		cloud = "aliyun"
	}
	switch cloud {
	case "general":
		go GenMail(to, from, title, body, mailType)
	case "aliyun":
		go AliMail(to, from, title, body, mailType)
	}
}
