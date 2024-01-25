// @Title
// @Description $
// @Author  55
// @Date  2022/5/12
package ipinfo

import (
	"github.com/zngw/zipinfo/config"
	"reflect"
	"strings"
)

var infoRegistry = make(map[string]Base)
var infos []Base

func registerInfo(info Base) {
	t := reflect.TypeOf(info).Elem()
	infoRegistry[strings.ToLower(t.Name())] = info
}

// 输入生效的第三方库及调用顺序
func Init(third []interface{}) {
	infos = nil
	if third != nil {
		for _, v := range third {
			if k, ok := v.(map[string]interface{})["name"]; ok {
				if info, ok := infoRegistry[k.(string)]; ok {
					if info.Init(v) {
						infos = append(infos, info)
					}
				}
			}
		}
	}

	return
}

func GetIpInfoFree(ip string, free bool) (err error, info *IpInfo) {
	if infos == nil {
		Init(config.Cfg.Third)
	}

	var result = InfoTypeFail
	for i, _ := range infos {
		tc := infos[i]
		if free && !tc.CanFree() {
			continue
		}

		tmp := tc.IpInfo(ip)
		if tmp == nil {
			continue
		}

		r := tmp.check()
		if r == InfoTypeFail {
			continue
		} else if r == InfoTypeSuccess {
			info = tmp
			result = r
			break
		} else {
			if result < r {
				info = tmp
				result = r
			}
		}
	}

	return
}

func GetIpInfo(ip string) (err error, info *IpInfo) {
	return GetIpInfoFree(ip, false)
}
