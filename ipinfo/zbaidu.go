// @Title
// @Description $
// @Author  55
// @Date  2022/5/15
package ipinfo

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"reflect"
	"strings"
	"time"
)

type BaiDu struct {
	Url  string
	Free bool
}

func init() {
	registerInfo(new(BaiDu))
}

func (s *BaiDu) Init(cfg interface{}) bool {
	m := cfg.(map[string]interface{})
	s.Url = m["url"].(string)
	s.Free = m["free"].(bool)

	return true
}

func (s *BaiDu) CanFree() bool {
	return s.Free
}

// 百度接口
// 无国家显示
func (s *BaiDu) IpInfo(ip string) (info *IpInfo) {
	info = new(IpInfo)
	info.Status = "fail"
	info.Type = reflect.TypeOf(s).Elem().Name()

	url := fmt.Sprintf(s.Url, ip)
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

	m := make(map[string]interface{})
	err = json.Unmarshal(body, &m)
	if err != nil {
		// 网页解析错误
		return
	}

	if m["status"] != "0" {
		return
	}

	var data interface{} = nil
	if d, ok := m["data"]; ok {
		arr := d.([]interface{})
		if len(arr) > 0 {
			data = arr[0]
		}
	}

	if data == nil {
		return
	}

	md := data.(map[string]interface{})
	if location, ok := md["location"]; ok {
		l := strings.Split(location.(string), " ")
		if len(l) != 2 {
			return
		}
		info.Isp = l[1]
		location = l[0]

		r := strings.Split(location.(string), "省")
		if len(r) == 2 {
			info.Region = strings.TrimSpace(r[0])
			c := strings.Split(r[1], "市")
			if len(c) >= 1 {
				info.City = strings.TrimSpace(c[0])
			}
			if len(c) == 2 {
				info.Isp = strings.TrimSpace(c[1])
			}
		} else if len(r) == 1 {
			// 自治区
			if strings.Index(r[0], "自治区") > 0 {
				c := strings.Split(r[0], "自治区")
				if len(c) == 2 {
					info.Region = strings.TrimSpace(c[0]) + "自治区"
					info.City = strings.TrimSpace(c[1])
				} else {
					return
				}
			} else {
				// 直辖市
				c := strings.Split(r[0], "市")
				if len(c) >= 2 {
					info.Region = strings.TrimSpace(c[0])
					info.City = strings.TrimSpace(c[1])
				}

				if len(c) == 3 {
					info.Isp = strings.TrimSpace(c[2])
				}

				if len(c) == 1 {
					info.Country = c[0]
				}
			}
		}
	}

	if info.Region == "" && info.City == "" {
		return
	}

	info.Status = "success"
	if info.isChinaRegion() {
		info.Country = "中国"
	} else {
		info.Country = "国外"
	}

	info.Query = ip

	return
}
