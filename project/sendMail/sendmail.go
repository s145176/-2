package sendMail

import (
	"fmt"
	"gopkg.in/gomail.v2"
)

func SendMail(body string) {
	user := "1451763509@qq.com"
	pass := "mfrpkarifocbhhad"
	host := "smtp.qq.com"
	port := 587
	m := gomail.NewMessage()
	m.SetHeader("From", user)
	//接收人
	m.SetHeader("To", user)

	m.SetHeader("Subject", "羽毛球")
	//内容
	m.SetBody("text/html", body)

	d := gomail.NewDialer(host, port, user, pass)
	if err := d.DialAndSend(m); err != nil {
		fmt.Println(err)
	}
}
