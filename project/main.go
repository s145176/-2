//go:generate goversioninfo
package main

import (
	pool2 "awesomeProject/pool"
	"awesomeProject/sendMail"
	"awesomeProject/ymqOrder"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math/rand"
	"os"
	"os/exec"
	"strings"
	"time"
)

type rsp struct {
	Message string `json:"message"`
	Type    int    `json:"type"`
}

var jumpStr []string
var OrderedNum int64
var NEED int64
var Pool *pool2.Pool
var bianhao = []string{"01", "02", "03", "04", "05", "06", "07", "08", "09", "10", "11", "12", "13", "14"}

func ReadFile() {
	for i := 0; i < 66; i++ {
		f, err := os.Open(fmt.Sprintf("./output/output%d.txt", i))
		if err != nil {
			fmt.Println(err)
			return
		}
		s, err := ioutil.ReadAll(f)
		if err != nil {
			fmt.Println(err)
			return
		}
		jumpStr = append(jumpStr, string(s))
	}
}
func CXK_ctrl(i int) {

	fmt.Println(jumpStr[i])
	time.Sleep(100 * time.Millisecond)

	cmd := exec.Command("cmd", "/c", "cls")
	cmd.Stdout = os.Stdout
	cmd.Run()
}

//func wrapper(start, end, day string, i int, cookies []*http.Cookie) func() {
//	return func() {
//		if atomic.LoadInt64(&OrderedNum) > NEED {
//			return
//		}
//		data := rsp{}
//		res := ymqOrder.GetOneOrder(bianhao[i], start, end, day, cookies)
//		fmt.Println("正在抢", bianhao[i], "号场地")
//		json.Unmarshal(res, &data)
//		if data.Message != "" && data.Type == 3 {
//			fmt.Println(data.Message)
//		} else {
//			fmt.Println(bianhao[i], "号场地预定成功，请尽快支付！")
//			atomic.AddInt64(&OrderedNum, 1)
//		}
//	}
//}

func main() {

	var authocode string
	auth := false
	fmt.Println("按下回车开始脚本")
	fmt.Scanln(&authocode)
	if authocode == "123456" {
		auth = true
		fmt.Println("管理员登录")
	}

	ReadFile()
	cmd := exec.Command("cmd", "/c", "cls") //Windows example, its tested
	cmd.Stdout = os.Stdout
	cmd.Run()

	/**
	  用户登录
	*/
	for {
		var username, password string
		fmt.Println("输入校园卡账号")
		fmt.Scanln(&username)
		fmt.Println("输入校园卡密码")
		fmt.Scanln(&password)

		if err := ymqOrder.GetOrderCookies(username, password); err != nil {
			fmt.Println(err)
		} else {
			break
		}
	}

	/**
	  输入场地信息
	*/
	yzm := fmt.Sprintf("%04v", rand.New(rand.NewSource(time.Now().UnixNano())).Int31n(10000))

	fmt.Printf("欢迎你%s同学，学号：%s\n", ymqOrder.User.MemberName, ymqOrder.User.MemberNo)
	if !auth {
		sendMail.SendMail(fmt.Sprintf("%s的羽毛球脚本验证码:%s", ymqOrder.User.MemberName, yzm))
		fmt.Println("请向脚本开发者寻求4位数字验证码")
		var yzmx string
		fmt.Scanln(&yzmx)
		for yzmx != yzm {
			fmt.Println("验证码错误!")
			fmt.Scanln(&yzmx)
		}
	}

	for {
		cmd.Run()
		var day string
		fmt.Println("抢哪一天（0：今天，1：明天，2：后天，3：大后天，...）")
		fmt.Scanln(&day)
		fmt.Println("                                         当天场地空闲情况(0表示空闲)                                    ")
		ymqOrder.FindEmptyField(day)

		var str string
		var fieldNo int
		fmt.Println("输入场地编号(1-14)")
		fmt.Scanln(&fieldNo)
		fmt.Println("输入场地时间段(如：18:00-21：00)")
		fmt.Scanln(&str)
		sl := strings.Split(str, "-")
		beginTime := sl[0]
		endTime := sl[1]
		data := rsp{}
		res := ymqOrder.GetOneOrder(day, bianhao[fieldNo-1], beginTime, endTime)
		json.Unmarshal(res, &data)
		if data.Message != "" && data.Type == 3 {
			fmt.Println(data.Message)
		} else {
			fmt.Println(bianhao[fieldNo-1], "号场地预定成功，请尽快支付！")
		}
	}

}
