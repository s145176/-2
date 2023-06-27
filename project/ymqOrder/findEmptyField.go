package ymqOrder

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"strconv"
	"strings"
)

/**
 * ┌───┐   ┌───┬───┬───┬───┐ ┌───┬───┬───┬───┐ ┌───┬───┬───┬───┐ ┌───┬───┬───┐
 * │Esc│   │ F1│ F2│ F3│ F4│ │ F5│ F6│ F7│ F8│ │ F9│F10│F11│F12│ │P/S│S L│P/B│  ┌┐    ┌┐    ┌┐
 * └───┘   └───┴───┴───┴───┘ └───┴───┴───┴───┘ └───┴───┴───┴───┘ └───┴───┴───┘  └┘    └┘    └┘
 * ┌───┬───┬───┬───┬───┬───┬───┬───┬───┬───┬───┬───┬───┬───────┐ ┌───┬───┬───┐ ┌───┬───┬───┬───┐
 * │~ `│! 1│@ 2│# 3│$ 4│% 5│^ 6│& 7│* 8│( 9│) 0│_ -│+ =│ BacSp │ │Ins│Hom│PUp│ │N L│ / │ * │ - │
 * ├───┴─┬─┴─┬─┴─┬─┴─┬─┴─┬─┴─┬─┴─┬─┴─┬─┴─┬─┴─┬─┴─┬─┴─┬─┴─┬─────┤ ├───┼───┼───┤ ├───┼───┼───┼───┤
 * │ Tab │ Q │ W │ E │ R │ T │ Y │ U │ I │ O │ P │{ [│} ]│ | \ │ │Del│End│PDn│ │ 7 │ 8 │ 9 │   │
 * ├─────┴┬──┴┬──┴┬──┴┬──┴┬──┴┬──┴┬──┴┬──┴┬──┴┬──┴┬──┴┬──┴─────┤ └───┴───┴───┘ ├───┼───┼───┤ + │
 * │ Caps │ A │ S │ D │ F │ G │ H │ J │ K │ L │: ;│" '│ Enter  │               │ 4 │ 5 │ 6 │   │
 * ├──────┴─┬─┴─┬─┴─┬─┴─┬─┴─┬─┴─┬─┴─┬─┴─┬─┴─┬─┴─┬─┴─┬─┴────────┤     ┌───┐     ├───┼───┼───┼───┤
 * │ Shift  │ Z │ X │ C │ V │ B │ N │ M │< ,│> .│? /│  Shift   │     │ ↑ │     │ 1 │ 2 │ 3 │   │
 * ├─────┬──┴─┬─┴──┬┴───┴───┴───┴───┴───┴──┬┴───┼───┴┬────┬────┤ ┌───┼───┼───┐ ├───┴───┼───┤ E││
 * │ Ctrl│    │Alt │         Space         │ Alt│    │    │Ctrl│ │ ← │ ↓ │ → │ │   0   │ . │←─┘│
 * └─────┴────┴────┴───────────────────────┴────┴────┴────┴────┘ └───┴───┴───┘ └───────┴───┴───┘
 */

func FindEmptyField(dateadd string) {
	fmt.Print("场地编号\\开始时间")
	for i := 8; i < 21; i++ {
		fmt.Printf("%4v:00", i)
	}
	fmt.Println()
	result := [14][]string{}

	for _, v := range []string{"0", "1", "2"} {
		r := findEmptyField(dateadd, v)
		for i := 0; i < 14; i++ {
			result[i] = append(result[i], r[i]...)
		}
	}
	m := map[string]string{"0": "0"}
	for i := 0; i < 14; i++ {
		if i+1 > 9 {
			fmt.Printf("          场地%d", i+1)

		} else {
			fmt.Printf("           场地%d", i+1)

		}
		for j := 0; j < 13; j++ {
			fmt.Printf("  %4v ", m[result[i][j]])
		}
		fmt.Println()
	}

}

func findEmptyField(dateadd, TimePeriod string) [][]string {
	//上午、中午、下午
	rawUrl := fmt.Sprintf("https://tybsouthgym.xidian.edu.cn/Field/GetVenueState?dateadd=%s&TimePeriod=%s&VenueNo=01&FieldTypeNo=001&_=1683956753648", dateadd, TimePeriod)

	url, _ := url.Parse(rawUrl)
	jar, _ := cookiejar.New(nil)
	jar.SetCookies(url, Cookies)
	c := http.Client{Jar: jar}
	req, _ := http.NewRequest("GET", rawUrl, nil)
	rsp, _ := c.Do(req)
	//fmt.Println(rsp, err)
	res, _ := ioutil.ReadAll(rsp.Body)
	m := make(map[string]string)
	json.Unmarshal(res, &m)
	obj := []map[string]string{}
	json.Unmarshal([]byte(m["resultdata"]), &obj)

	x := 8
	y := 3
	if TimePeriod == "0" {
		x = 8
		y = 4
	} else if TimePeriod == "1" {
		x = 12
		y = 6
	} else {
		x = 18
		y = 3
	}
	empty := make([][]string, 14)
	for i := range empty {
		empty[i] = make([]string, y)
	}

	for _, v := range obj {
		_, fieldNo, _ := strings.Cut(v["FieldName"], "羽毛球")
		n, _ := strconv.Atoi(fieldNo)
		bgt, _, _ := strings.Cut(v["BeginTime"], ":")
		b, _ := strconv.Atoi(bgt)
		empty[n-1][b-x] = v["FieldState"]
	}

	return empty
}
