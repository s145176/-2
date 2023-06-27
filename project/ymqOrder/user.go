package ymqOrder

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/http/cookiejar"
	url2 "net/url"
)

var User user

type user struct {
	MemberNo   string `json:"memberno"`
	MemberName string `json:"membername"`
}

func getMemberInfo() {
	url := "https://tybsouthgym.xidian.edu.cn/User/GetMemberInfo?_=1653641085587"
	req, _ := http.NewRequest("GET", url, nil)
	parseUrl, _ := url2.Parse(url)
	jar, _ := cookiejar.New(nil)
	jar.SetCookies(parseUrl, Cookies)
	client := http.Client{Jar: jar}
	rsp, _ := client.Do(req)
	data, _ := ioutil.ReadAll(rsp.Body)
	json.Unmarshal(data[1:len(data)-1], &User)
}
