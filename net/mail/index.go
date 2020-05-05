package mail

import "github.com/xurwxj/viper"

// need following config in config.json:
// "email": {
// 		"cloud": "aliyun",
// 		"account": "xxx@xxx.com",
// 		"pwd": "xxx",
// 		"host": "smtpdm.xxx.com:25",
// 		"service": {
// 			"panic": "xxx@xurw.com"
// 		}
//   },
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
