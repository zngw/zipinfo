// @Title
// @Description $
// @Author  55
// @Date  2022/5/12
package ipinfo

const (
	InfoTypeFail      = 0 // 获取失败
	InfoTypeNotRegion = 1 // 无省份信息
	InfoTypeNotCity   = 2 // 无城市信息
	InfoTypeSuccess   = 3 // 获取成功
)

type Base interface {
	Init(cfg interface{}) bool
	IpInfo(ip string) *IpInfo
	CanFree() bool
}
