// @Title
// @Description $
// @Author  55
// @Date  2022/5/15
package ipinfo

import (
	"bytes"
	"encoding/json"
	"fmt"
	"golang.org/x/text/encoding/simplifiedchinese"
	"golang.org/x/text/transform"
	"io"
	"io/ioutil"
	"net/http"
	"reflect"
	"strings"
	"time"
)

type PcOnline struct {
	Url  string
	Free bool
}

func init() {
	registerInfo(new(PcOnline))
}

func (s *PcOnline) Init(cfg interface{}) bool {
	m := cfg.(map[string]interface{})
	s.Url = m["url"].(string)
	s.Free = m["free"].(bool)

	return true
}

func (s *PcOnline) CanFree() bool {
	return s.Free
}

// 太平洋电脑网接口
// 无国家显示
func (s *PcOnline) IpInfo(ip string) (info *IpInfo) {
	info = new(IpInfo)
	info.Status = "fail"
	info.Type = reflect.TypeOf(s).Elem().Name()

	url := fmt.Sprintf(s.Url, ip)
	client := &http.Client{
		Timeout: time.Millisecond * 3000,
	}
	req, _ := http.NewRequest("GET", url, nil)
	resp, err := client.Do(req)
	if err != nil {
		// 读取网页数据错误
		return
	}

	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil || resp.StatusCode != 200 {
		// 读取网页数据错误
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

	m := make(map[string]interface{})
	err = json.Unmarshal(body, &m)
	if err != nil {
		// 网页解析错误
		return
	}

	info.Status = "success"
	if r, ok := m["pro"]; ok {
		info.Region = strings.ReplaceAll(strings.ReplaceAll(r.(string), "省", ""), "市", "")
	}

	if city, ok := m["city"]; ok {
		info.City = strings.ReplaceAll(city.(string), "市", "")
	}
	if info.isChinaRegion() {
		info.Country = "中国"
	} else {
		info.Country = "国外"
	}
	if isp, ok := m["addr"]; ok {
		i := strings.Split(isp.(string), " ")
		if len(i) == 2 {
			info.Isp = i[1]
		} else {
			info.Isp = isp.(string)
		}
	}

	info.Query = ip

	return
}
