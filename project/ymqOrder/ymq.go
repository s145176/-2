package ymqOrder

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/cookiejar"
	url2 "net/url"
)

type PostResult struct {
	E int    `json:"e"`
	M string `json:"m"`
}

func GetOneOrder(day, fieldno, begin, end string) []byte {

	param := url2.Values{}
	s := fmt.Sprintf("[{\"FieldNo\":\"%s\",\"FieldTypeNo\":\"%s\",\"FieldName\":\"%s\",\"BeginTime\":\"%s\",\"Endtime\":\"%s\",\"Price\":\"%s\"}]",
		"YMQ0"+fieldno,
		"001",
		"羽毛球"+fieldno,
		begin,
		end,
		"6.00")

	param.Set("checkdata", s)
	param.Set("dateadd", day) //日期增量
	param.Set("VenueNo", "01")
	urlStr := "https://tybsouthgym.xidian.edu.cn/Field/OrderField?" + param.Encode()
	r, _ := http.NewRequest("GET", urlStr, nil)
	parseUrl, _ := url2.Parse(urlStr)

	jar, _ := cookiejar.New(nil)
	jar.SetCookies(parseUrl, Cookies)
	client := &http.Client{Jar: jar}
	response, err := client.Do(r)
	if err != nil {
		fmt.Println(err)
	}
	data, _ := ioutil.ReadAll(response.Body)
	return data
}
