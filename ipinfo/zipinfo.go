// @Title
// @Description $
// @Author  55
// @Date  2022/6/23
package ipinfo

import (
	"encoding/json"
	"fmt"
	"github.com/zngw/zipinfo/util"
	"io"
	"net/http"
	"reflect"
	"time"
)

type ZipInfo struct {
	Url   string
	Free  bool
	User  string
	Token string
}

func init() {
	registerInfo(new(ZipInfo))
}

func (s *ZipInfo) Init(cfg interface{}) bool {
	m := cfg.(map[string]interface{})
	s.Url = m["url"].(string)
	s.Free = m["free"].(bool)
	s.User = m["user"].(string)
	s.Token = m["token"].(string)

	return true
}

func (s *ZipInfo) CanFree() bool {
	return s.Free
}

func (s *ZipInfo) IpInfo(ip string) (info *IpInfo) {
	info = new(IpInfo)
	info.Status = "fail"
	info.Type = reflect.TypeOf(s).Elem().Name()

	url := fmt.Sprintf("%s/?ip=%s&user=%s&key=%s", s.Url, ip, s.User, util.Md5Str(ip+s.User+s.Token))
	client := &http.Client{
		Timeout: time.Millisecond * 1000,
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

	err = json.Unmarshal(body, &info)
	if err != nil {
		return
	}

	return
}
