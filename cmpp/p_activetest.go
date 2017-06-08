/*p_activetest.go
链路检测（CMPP_ACTIVE_TEST）操作
本操作仅适用于通信双方采用长连接通信方式时用于保持连接。
*/
package cmpp

import (
	"time"
)

//检测发包时间间隔 C=3分钟
var Interval time.Duration

//等回复时间 T=60秒
var WaitTime time.Duration

//连续发包限制最大次数 N=3
var CheckLimitNum int
var CheckCount int

//最后交互时间
var LastActive time.Time

type Active_TestBody struct{}

type Active_TestBody_Resp struct {
	Reserved uint8 //1一字节任意内容   0x00
}

func Up_LastAction() {
	LastActive = time.Now()
}
func Check() (IsActive bool) {
	IsActive = true
	NowTime := time.Now()
	SubTime := NowTime.Sub(LastActive)
	if SubTime >= Interval*time.Second {
		//发检测包
		IsActive = false
	}
	return IsActive
}
func Creat_Active_TestBody_Resp() (resp *Active_TestBody_Resp) {
	resp = new(Active_TestBody_Resp)
	resp.Reserved = 0x00
	return resp
}
func (resp *Active_TestBody_Resp) Encode_Active_TestBody_Resp() (body_buff []byte) {
	body_buff = make([]byte, 1)
	body_buff = append(body_buff, 0x00)
	return body_buff
}
