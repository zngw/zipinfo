// @Title
// @Description $
// @Author  55
// @Date  2022/5/12
package httputil

import (
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"github.com/zngw/log"
	"io/ioutil"
	"net/url"
	"sort"
)

var users = make(map[string]string)

// 初始化验证用户数据
func init() {
	raw, err := ioutil.ReadFile("user.json")
	if err != nil {
		log.Error("sys", "无本用户配置文件，以无验证模式启动")
		return
	}

	err = json.Unmarshal(raw, &users)
	if err != nil {
		log.Error("sys", "解析基本配置文件失败,以无验证模式启动")
		return
	}
}

// 验证身份
func Authentication(param url.Values) bool {
	var auth string
	var keys []string
	for k, _ := range param {
		if k != "key" {
			keys = append(keys, k)

			if k == "user" {
				auth = param.Get(k)
			}
		}
	}
	sort.Strings(keys)

	Key, ok := getUserKey(auth)
	if !ok {
		return false
	}

	if len(Key) == 0 {
		return true
	}

	data := ""
	for _, k := range keys {
		data += param.Get(k)
	}
	data += Key

	return param.Get("key") == md5StrToHex(data)
}

func md5StrToHex(txt string) string {
	m := md5.New()
	m.Write([]byte(txt))
	return hex.EncodeToString(m.Sum(nil))
}

func getUserKey(user string) (key string, ok bool) {
	if len(users) == 0 {
		ok = true
		return
	}

	key, ok = users[user]
	return
}
