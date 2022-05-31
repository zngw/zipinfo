// @Title
// @Description $
// @Author  55
// @Date  2022/5/16
package ipinfo

import "strings"

type IpInfo struct {
	Status  string `json:"status,omitempty"`
	Type    string `json:"type,omitempty"`
	Country string `json:"country,omitempty"`
	Region  string `json:"region,omitempty"`
	City    string `json:"city,omitempty"`
	Isp     string `json:"isp,omitempty"`
	Query   string `json:"query,omitempty"`
	Param   string `json:"param,omitempty"`
}

var region = "北京,天津,河北,山西,内蒙古,辽宁,吉林,黑龙江,上海,江苏," +
	"浙江,安徽,福建,江西,山东,河南,湖北,湖南,广东,广西,海南,重庆,四川," +
	"贵州,云南,西藏,陕西,甘肃,青海,宁夏,新疆,香港,澳门,台湾," +
	"广西壮族自治区,内蒙古自治区,新疆维吾尔自治区,宁夏回族自治区,西藏自治区"

//
func (info *IpInfo) isChinaRegion() bool {
	if info.Region == "" {
		return false
	}
	return strings.Index(region, info.Region) >= 0
}

func (info *IpInfo) isChina() bool {
	return "中国" == info.Country
}

func (info *IpInfo) check() (result int) {
	// 状态为失败
	if info.Status == "fail" {
		return InfoTypeFail
	}

	// 国外ip不涉及省市，直接返回成功
	if !info.isChina() {
		return InfoTypeSuccess
	}

	// 判断省
	if info.Region == "" || !info.isChinaRegion() {
		return InfoTypeNotRegion
	}

	if info.City == "" || info.City == "XX" {
		return InfoTypeNotCity
	}

	return InfoTypeSuccess
}
