// @Title
// @Description $
// @Author  55
// @Date  2022/5/15
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

type BaiDu struct {
}

func init() {
	registerInfo(new(BaiDu))
}

// 太平洋电脑网接口
// 无国家显示
func (s *BaiDu) IpInfo(ip string) (info *IpInfo) {
	info = new(IpInfo)
	info.Status = "fail"
	info.Type = reflect.TypeOf(s).Elem().Name()

	url := fmt.Sprintf("https://opendata.baidu.com/api.php?query=%s&co=&resource_id=6006&oe=utf8", ip)
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
