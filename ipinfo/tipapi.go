// @Title
// @Description $
// @Author  55
// @Date  2022/5/12
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

type IpApi struct {
}

func init() {
	registerInfo(new(IpApi))
}

// 通过ip-api.com获取ip地理位置信息
// 由于ip-api.com是国外的网站，对国内市级ip位置有一定误差
func (s *IpApi) IpInfo(ip string) (info *IpInfo) {
	info = new(IpInfo)
	info.Status = "fail"
	info.Type = reflect.TypeOf(s).Elem().Name()

	url := fmt.Sprintf("http://ip-api.com/json/%s?lang=zh-CN", ip)
	client := &http.Client{
		Timeout: time.Millisecond * 500,
	}
	req, _ := http.NewRequest("GET", url, nil)
	resp, err := client.Do(req)
	if err != nil {
		// 获取不到地理位置，
		return
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		// 读取网页数据错误
		return
	}
	if resp.StatusCode != 200 {
		return
	}
	m := make(map[string]interface{})
	err = json.Unmarshal(body, &m)
	if err != nil {
		// 网页解析错误
		return
	}

	info.Status = "success"
	if c, ok := m["country"]; ok {
		info.Country = c.(string)
	}

	if r, ok := m["regionName"]; ok {
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
