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
	"time"
)

type TaoBao struct {
}

func init() {
	registerInfo(new(TaoBao))
}

// 淘宝ip地址库
// 限制频率:每个用户的访问频率需大于1qps
func (s *TaoBao) IpInfo(ip string) (info *IpInfo) {
	info = new(IpInfo)
	info.Status = "fail"
	info.Type = reflect.TypeOf(s).Elem().Name()

	url := fmt.Sprintf("https://ip.taobao.com/outGetIpInfo?ip=%s&accessKey=alibaba-inc", ip)
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

	var result map[string]interface{}
	err = json.Unmarshal(body, &result)
	if err != nil {
		// 网页解析错误
		return
	}

	v, ok := result["code"]
	if !ok || int(v.(float64)) != 0 {
		// 数据错误
		return
	}

	d := result["data"]
	data, err := json.Marshal(d)
	if err != nil {
		return
	}

	err = json.Unmarshal(data, &info)
	if err != nil {
		return
	}

	info.Status = "success"
	info.Query = ip

	return
}
