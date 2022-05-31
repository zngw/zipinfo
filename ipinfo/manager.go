// @Title
// @Description $
// @Author  55
// @Date  2022/5/12
package ipinfo

import (
	"reflect"
)

var defaultEnable = []string{"TaoBao", "UserAgentInfo", "Net126", "BaiDu", "PcOnline", "IpApi"}
var infoRegistry = make(map[string]Base)
var infos []Base

func registerInfo(info Base) {
	t := reflect.TypeOf(info).Elem()
	infoRegistry[t.Name()] = info

	Init(defaultEnable)
}

// 输入生效的第三方库及调用顺序
func Init(enables []string) {
	if enables != nil {
		for _, t := range enables {
			if info, ok := infoRegistry[t]; ok {
				infos = append(infos, info)
			}
		}
	}

	return
}

func GetIpInfo(ip string) (err error, info *IpInfo) {
	var result = InfoTypeFail
	for i, _ := range infos {
		tc := infos[i]
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
