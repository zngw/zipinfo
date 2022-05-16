// @Title
// @Description $
// @Author  55
// @Date  2022/5/14
package ipinfo

import (
	"bytes"
	"fmt"
	"golang.org/x/text/encoding/simplifiedchinese"
	"golang.org/x/text/transform"
	"io/ioutil"
	"net/http"
	"reflect"
	"regexp"
	"strings"
	"time"
)

type Net126 struct {
}

func init() {
	registerInfo(new(Net126))
}

// 通过网易126.net获取ip地理位置信息
// 无国家显示
func (s *Net126) IpInfo(ip string) (info *IpInfo) {
	info = new(IpInfo)
	info.Status = "fail"
	info.Type = reflect.TypeOf(s).Elem().Name()

	url := fmt.Sprintf("https://ip.ws.126.net/ipquery?ip=%s", ip)
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

	contentType := strings.ToLower(resp.Header.Get("Content-Type"))
	if strings.Contains(contentType, "gbk") {
		reader := transform.NewReader(bytes.NewReader(body), simplifiedchinese.GBK.NewDecoder())
		body, err = ioutil.ReadAll(reader)
		if err != nil {
			return
		}
	}

	data := string(body)

	pattern := `{city:\"(.*)\", province:\"(.*)\"}`
	myRegex, err := regexp.Compile(pattern)
	if err != nil {
		return
	}

	params := myRegex.FindStringSubmatch(data)
	if len(params) != 3 {
		return
	}

	info.Status = "success"
	info.City = strings.ReplaceAll(params[1], "市", "")
	info.Region = strings.ReplaceAll(strings.ReplaceAll(params[2], "省", ""), "市", "")
	if info.isChinaRegion() {
		info.Country = "中国"
	} else {
		info.Country = "国外"
	}
	info.Query = ip

	return
}
