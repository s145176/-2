package ymqOrder

import (
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/cookiejar"
	url2 "net/url"
	"strings"
)

var Cookies []*http.Cookie

func GetOrderCookies(username, password string) error {
	//登录url
	postUrl := "https://xxcapp.xidian.edu.cn/uc/wap/login/check"

	userValue := url2.Values{}
	userValue.Set("username", username)
	userValue.Set("password", password)

	loginReq, _ := http.NewRequest(http.MethodPost, postUrl, strings.NewReader(userValue.Encode()))
	loginReq.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	client := &http.Client{Transport: &http.Transport{TLSClientConfig: &tls.Config{InsecureSkipVerify: true}}}

	rsp, _ := client.Do(loginReq)
	data, _ := ioutil.ReadAll(rsp.Body)
	postResult := &PostResult{}
	json.Unmarshal(data, postResult)
	if postResult.E == 0 {
		fmt.Println("waiting...")
	} else {
		return fmt.Errorf(postResult.M + "请重新尝试登录")
	}

	indexUrl := "https://xxcapp.xidian.edu.cn//uc/api/oauth/index?appid=200201218103247434&redirect=http%3a%2f%2ftybsouthgym.xidian.edu.cn%2fUser%2fQYLogin&state=STATE"
	parseIndexUrl, _ := url2.Parse(indexUrl)
	indexReq, _ := http.NewRequest("GET", indexUrl, nil)
	jar1, _ := cookiejar.New(nil)
	jar1.SetCookies(parseIndexUrl, rsp.Cookies())
	c := &http.Client{CheckRedirect: func(req *http.Request, via []*http.Request) error {
		if len(via) == 3 {
			return http.ErrUseLastResponse
		}
		return nil
	}, Jar: jar1}
	rsp, _ = c.Do(indexReq)

	indexCookies := rsp.Cookies()
	Cookies = indexCookies
	getMemberInfo()
	return nil
}
