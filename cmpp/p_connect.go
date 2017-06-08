/*p_connect.go
SP请求连接到ISMG（CMPP_CONNECT）操作
CMPP_CONNECT操作的目的是SP向ISMG注册作为一个合法SP身份，若注册成功后即建立了应用层的连接，
此后SP可以通过此ISMG接收和发送短信。
ISMG以CMPP_CONNECT_RESP消息响应SP的请求。
*/
package cmpp

import (
	"../comm"
	"errors"
	//"fmt"
)

//即SP的企业代码
var Source_Addr string

//秘匙
var Shared_secret string

//version 3.0
const Sys_version = uint8(0x30)

type ConnectBody struct {
	Source_Addr string //6字节
	/*源地址，此处为SP_Id，即SP的企业代码。*/

	AuthenticatorSource string //16字节
	/*用于鉴别源地址。其值通过单向MD5 hash计算得出，表示如下：
	AuthenticatorSource = MD5（Source_Addr+9 字节的0 +shared secret+timestamp）
	Shared secret 由中国移动与源地址实体事先商定，timestamp格式为：MMDDHHMMSS，即月日时分秒，10位。*/

	Version uint8 //1字节
	/*双方协商的版本号(高位4bit表示主版本号,低位4bit表示次版本号)，
	对于3.0的版本，高4bit为3，低4位为0*/

	Timestamp uint32 //4字节
	/*时间戳的明文,由客户端产生,格式为MMDDHHMMSS，即月日时分秒，10位数字的整型，右对齐。*/
}
type ConnectBody_Resp struct {
	Status uint32 //4字节
	/*状态
	0：正确
	1：消息结构错
	2：非法源地址
	3：认证错
	4：版本太高
	5~ ：其他错误*/

	AuthenticatorISMG string //16字节
	/*ISMG认证码，用于鉴别ISMG。其值通过单向MD5 hash计算得出，表示如下：
	AuthenticatorISMG =MD5（Status+AuthenticatorSource+shared secret），Shared secret 由中国移动与源地址实体事先商定，AuthenticatorSource为源地址实体发送给ISMG的对应消息CMPP_Connect中的值。
	认证出错时，此项为空。*/

	Version uint8 //1字节
	/*服务器支持的最高版本号，对于3.0的版本，高4bit为3，低4位为0*/
}

/*SP或ISMG请求拆除连接（CMPP_TERMINATE）操作
CMPP_TERMINATE操作的目的是SP或ISMG基于某些原因决定拆除当前的应用层连接而发起的操作。
此操作完成后SP与ISMG之间的应用层连接被释放，此后SP若再要与ISMG通信时应发起CMPP_CONNECT操作。
ISMG或SP以CMPP_TERMINATE_RESP消息响应请求。
无消息体。*/
type TerminateBody struct{}
type TerminateBody_Resp struct{}

func Creat_ConnectBody() (body *ConnectBody) {

	body = new(ConnectBody)
	body.Source_Addr = Source_Addr
	md5_buff := Creat_AuthenticatorSource()
	body.AuthenticatorSource = string(md5_buff)
	body.Version = Sys_version
	body.Timestamp = comm.Get_Timestamp_uint32()

	return body

}
func (this *ConnectBody) Encode_ConnectBody() (body_buff []byte, totall_len int) {

	Source_Addr_byte := comm.Autocompletion([]byte(this.Source_Addr), 6)
	Auth_byte := comm.Autocompletion([]byte(this.AuthenticatorSource), 16)
	Version_byte := comm.Uint8_byte(this.Version)
	Timestamp_byte := comm.Uint32_byte(this.Timestamp)

	totall_len = 27
	body_buff = make([]byte, totall_len)
	copy(body_buff[0:], Source_Addr_byte[0:6])
	copy(body_buff[6:], Auth_byte[0:16])
	copy(body_buff[22:], Version_byte[0:1])
	copy(body_buff[23:], Timestamp_byte[0:4])
	return body_buff, totall_len
}
func Decode_ConnectBody_Resp(body_buff []byte) (resp *ConnectBody_Resp, err error) {

	if len(body_buff) < 21 {
		return nil, errors.New("Decode_ConnectBody_Resp:Buffer len no Enough 21")
	}
	resp = new(ConnectBody_Resp)
	resp.Status, err = comm.Byte_uint32(body_buff[0:4])
	if err != nil {
		return nil, err
	}
	resp.AuthenticatorISMG = string(body_buff[4:20])
	resp.Version, _ = comm.Byte_uint8(body_buff[20:21])
	return resp, nil

}

func Creat_AuthenticatorSource() (md5 []byte) {
	//AuthenticatorSource = MD5（Source_Addr+9 字节的0 +shared secret+timestamp）

	Timestamp_str := comm.Get_now("mdHis")
	Source_buff := []byte(Source_Addr)
	secret_buff := []byte(Shared_secret)
	Timestamp_buff := []byte(Timestamp_str)
	lens := len(Source_buff) + 9 + len(secret_buff) + len(Timestamp_buff)
	Auth_buff := make([]byte, lens)
	copy(Auth_buff[0:], Source_buff[0:])
	copy(Auth_buff[len(Source_buff)+9:], secret_buff[0:])
	copy(Auth_buff[len(Source_buff)+9+len(secret_buff):], Timestamp_buff[0:])
	//fmt.Println(Auth_buff)
	md5_buff := comm.Md5_encode(Auth_buff)
	return md5_buff

}
