package mail

import "github.com/xurwxj/viper"

func SendMail(to, from, title, body, mailType string) error {
	cloud := viper.GetString("email.cloud")
	if cloud == "" {
		cloud = "aliyun"
	}
	switch cloud {
	case "general":
		return GenMail(to, from, title, body, mailType)
	case "aliyun":
		return AliMail(to, from, title, body, mailType)
	}
	return nil
}
