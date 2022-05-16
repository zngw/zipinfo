// @Title
// @Description $
// @Author  55
// @Date  2022/5/16
package ipinfo

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"reflect"
	"strings"
	"time"
)

type UserAgentInfo struct {
}

func init() {
	registerInfo(new(UserAgentInfo))
}

// ip.useragentinfo.com
// 直接调用即可【没有频率限制】
func (s *UserAgentInfo) IpInfo(ip string) (info *IpInfo) {
	info = new(IpInfo)
	info.Status = "fail"
	info.Type = reflect.TypeOf(s).Elem().Name()

	url := fmt.Sprintf("https://ip.useragentinfo.com/json?ip=%s", ip)
	client := &http.Client{
		Timeout: time.Millisecond * 500,
	}
	req, _ := http.NewRequest("GET", url, nil)
	resp, err := client.Do(req)
	if err != nil {
		// 读取网页数据错误
		return
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil || resp.StatusCode != 200 {
		// 读取网页数据错误
		return
	}

	var m map[string]interface{}
	err = json.Unmarshal(body, &m)
	if err != nil {
		// 网页解析错误
		return
	}

	info.Status = "success"
	if c, ok := m["country"]; ok {
		info.Country = c.(string)
	}

	if r, ok := m["province"]; ok {
		info.Region = strings.ReplaceAll(strings.ReplaceAll(r.(string), "省", ""), "市", "")
	}

	if city, ok := m["city"]; ok {
		info.City = strings.ReplaceAll(city.(string), "市", "")
	}

	if i, ok := m["isp"]; ok {
		info.Isp = i.(string)
	}

	info.Query = ip

	return
}
