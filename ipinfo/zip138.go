// @Title
// @Description $
// @Author  55
// @Date  2022/6/23
package ipinfo

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"reflect"
	"time"
)

type Ip138 struct {
	Url   string
	Free  bool
	Token string
}

func init() {
	registerInfo(new(Ip138))
}

func (s *Ip138) Init(cfg interface{}) bool {
	m := cfg.(map[string]interface{})
	s.Url = m["url"].(string)
	s.Free = m["free"].(bool)
	s.Token = m["token"].(string)

	return true
}

func (s *Ip138) CanFree() bool {
	return s.Free
}

func (s *Ip138) IpInfo(ip string) (info *IpInfo) {
	info = new(IpInfo)
	info.Status = "fail"
	info.Type = reflect.TypeOf(s).Elem().Name()

	url := fmt.Sprintf(s.Url, ip)
	client := &http.Client{
		Timeout: time.Millisecond * 500,
	}
	req, _ := http.NewRequest("GET", url, nil)
	req.Header.Add("token", s.Token)
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

	var jsonInfo struct {
		Ret  string    `json:"ret"`
		Ip   string    `json:"ip"`
		Data [6]string `json:"data"`
	}

	err = json.Unmarshal(body, &jsonInfo)
	if err != nil {
		return
	}

	info.Status = "success"
	info.Country = jsonInfo.Data[0]
	info.Region = jsonInfo.Data[1]
	info.City = jsonInfo.Data[2]
	info.Isp = jsonInfo.Data[4]
	info.Query = ip

	return
}
