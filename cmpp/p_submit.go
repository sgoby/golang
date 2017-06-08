/*p_submit.go
SP向ISMG提交短信（CMPP_SUBMIT）操作
CMPP_SUBMIT操作的目的是SP在与ISMG建立应用层连接后向ISMG提交短信。
ISMG以CMPP_SUBMIT_RESP消息响应。
*/
package cmpp

import (
	"../comm"
	"errors"
	"time"
)

type SubmitBody struct {
	Msg_Id uint64 //8
	/*信息标识。*/

	Pk_total uint8 //1
	/*相同Msg_Id的信息总条数，从1开始。*/

	Pk_number uint8 //1
	/*相同Msg_Id的信息序号，从1开始。*/

	Registered_Delivery uint8 //1
	/*是否要求返回状态确认报告：
	  0：不需要；
	  1：需要。*/

	Msg_level uint8 //1
	/*信息级别。*/

	Service_Id string //10
	/*业务标识，是数字、字母和符号的组合。*/

	Fee_UserType uint8 //1
	/*计费用户类型字段：
	  0：对目的终端MSISDN计费；
	  1：对源终端MSISDN计费；
	  2：对SP计费；
	  3：表示本字段无效，对谁计费参见Fee_terminal_Id字段。*/

	Fee_terminal_Id string //32
	/*被计费用户的号码，当Fee_UserType为3时该值有效，当Fee_UserType为0、1、2时该值无意义。*/

	Fee_terminal_type uint8 //1
	/* 被计费用户的号码类型，0：真实号码；1：伪码。*/

	TP_pId uint8 //1
	/*GSM协议类型。详细是解释请参考GSM03.40中的9.2.3.9。 默认为0*/

	TP_udhi uint8 //1
	/*GSM协议类型。详细是解释请参考GSM03.40中的9.2.3.23,仅使用1位，右对齐。默认为0*/

	Msg_Fmt uint8 //1
	/*信息格式：
	  0：ASCII串；
	  3：短信写卡操作；
	  4：二进制信息；
	  8：UCS2编码；
	  15：含GB汉字。。。。。。*/

	Msg_src string //6
	/* 信息内容来源(SP_Id)。*/

	FeeType string //2
	/*资费类别：
	  01：对“计费用户号码”免费；
	  02：对“计费用户号码”按条计信息费；
	  03：对“计费用户号码”按包月收取信息费。*/

	FeeCode string //6
	/*资费代码（以分为单位）。*/

	ValId_Time string //17
	/*存活有效期，格式遵循SMPP3.3协议。 默认为空*/

	At_Time string //17
	/*定时发送时间，格式遵循SMPP3.3协议。默认为空*/

	Src_Id string //21
	/*源号码。SP的服务代码或前缀为服务代码的长号码, 网关将该号码完整的填到SMPP协议
	  Submit_SM消息相应的source_addr字段，该号码最终在用户手机上显示为短消息的主叫号码。*/

	DestUsr_tl uint8 //1
	/*接收信息的用户数量(小于100个用户)。*/

	Dest_terminal_Id string //32*DestUsr_tl
	/*接收短信的MSISDN号码。*/

	Dest_terminal_type uint8 //1
	/*接收短信的用户的号码类型，0：真实号码；1：伪码。*/

	Msg_Length uint8 //1
	/*信息长度(Msg_Fmt值为0时：<160个字节；其它<=140个字节)，取值大于或等于0。*/

	Msg_Content []byte //Msg_length string
	/*信息内容。*/

	LinkID string //20
	/*点播业务使用的LinkID，非点播类业务的MT流程不使用该字段。*/

}

/*
提交状态
*/
type SubmitStatus struct {
	Msg_Id      uint64
	Result      uint32
	LinkID      string
	Sequence_Id uint32
}

//
var SubmitStatus_map map[uint32]*SubmitStatus

type SubmitBody_Resp struct {
	Msg_Id uint64 //8
	/*信息标识，生成算法如下：
	  采用64位（8字节）的整数：
	  （1）时间（格式为MMDDHHMMSS，即月日时分秒）：bit64~bit39，其中
	  bit64~bit61：月份的二进制表示；
	  bit60~bit56：日的二进制表示；
	  bit55~bit51：小时的二进制表示；
	  bit50~bit45：分的二进制表示；
	  bit44~bit39：秒的二进制表示；
	  （2）短信网关代码：bit38~bit17，把短信网关的代码转换为整数填写到该字段中；
	  （3）序列号：bit16~bit1，顺序增加，步长为1，循环使用。
	  各部分如不能填满，左补零，右对齐。
	  （SP根据请求和应答消息的Sequence_Id一致性就可得到CMPP_Submit消息的Msg_Id）*/
	Result uint32 //4
	/*结果：
	  0：正确；
	  1：消息结构错；
	  2：命令字错；
	  3：消息序号重复；
	  4：消息长度错；
	  5：资费代码错；
	  6：超过最大信息长；
	  7：业务代码错；
	  8：流量控制错；
	  9：本网关不负责服务此计费号码；
	  10：Src_Id错误；
	  11：Msg_src错误；
	  12：Fee_terminal_Id错误；
	  13：Dest_terminal_Id错误；*/
}

//msgid 中的序号
var Submit_count uint16

//网关代码
var IsmgID uint64

func Creat_SubmitBody() (body *SubmitBody) {

	body = new(SubmitBody)

	body.Msg_Id = 0              // uint64 //8  Creat_MsgId()
	body.Pk_total = 1            // uint8 //1
	body.Pk_number = 1           // uint8 //1
	body.Registered_Delivery = 1 //uint8 //1
	body.Msg_level = 1           // uint8 //1
	body.Service_Id = ""         // #string //10 业务代码
	body.Fee_UserType = 0        // uint8 //1 对目的终端计费
	body.Fee_terminal_Id = ""    // string //32  用户手号码
	body.Fee_terminal_type = 0   // uint8 //1
	body.TP_pId = 0              // uint8 //1
	body.TP_udhi = 0             // uint8 //1
	body.Msg_Fmt = 8             // uint8 //1
	body.Msg_src = Source_Addr   // string //6
	body.FeeType = "02"          // string //2
	body.FeeCode = "100"         // #string //6
	body.ValId_Time = ""         // string //17
	body.At_Time = ""            // string //17
	body.Src_Id = ""             // #string //21  长号码
	body.DestUsr_tl = 1          // uint8 //1
	body.Dest_terminal_Id = ""   // #string //32*DestUsr_tl 接收短信的MSISDN号码  即手机86手机号码
	body.Dest_terminal_type = 0  // uint8 //1
	body.Msg_Length = 0          // #uint8 //1
	//body.Msg_Content = &[]byte{} // #string //Msg_length string
	body.LinkID = "" //# string //20

	return body
}
func Creat_MsgId() uint64 {

	now := time.Now()
	mm := uint64(now.Month())
	dd := uint64(now.Day())
	hh := uint64(now.Hour())
	ii := uint64(now.Minute())
	ss := uint64(now.Second())

	step := uint64(Submit_count)
	timestamp := uint64(0)
	//时间（格式为MMDDHHMMSS，即月日时分秒）：bit64~bit39
	timestamp = timestamp | (mm << 60) | (dd << 55) | (hh << 50) | (ii << 44) | (ss << 38)
	//短信网关代码
	timestamp = timestamp | (IsmgID << 16)
	//序列号：bit16~bit1，顺序增加，步长为1，循环使用。
	timestamp = timestamp | step

	return timestamp
}
func (this *SubmitBody) Encode_SubmitBody() (body_buff []byte, totall_len int) {

	Body_len := 0
	temp_buff := make([]byte, 500)

	Msg_Id_Buff := comm.Uint64_byte(this.Msg_Id) //uint64 //8
	copy(temp_buff[Body_len:], Msg_Id_Buff[0:])
	Body_len += len(Msg_Id_Buff)

	Pk_total_Buff := comm.Uint8_byte(this.Pk_total) //uint8 //1
	copy(temp_buff[Body_len:], Pk_total_Buff[0:])
	Body_len += len(Pk_total_Buff)

	Pk_number_Buff := comm.Uint8_byte(this.Pk_number) //uint8 //1
	copy(temp_buff[Body_len:], Pk_number_Buff[0:])
	Body_len += len(Pk_number_Buff)

	Registered_Delivery_Buff := comm.Uint8_byte(this.Registered_Delivery) //uint8 //1
	copy(temp_buff[Body_len:], Registered_Delivery_Buff[0:])
	Body_len += len(Registered_Delivery_Buff)

	Msg_level_Buff := comm.Uint8_byte(this.Msg_level) //uint8 //1
	copy(temp_buff[Body_len:], Msg_level_Buff[0:])
	Body_len += len(Msg_level_Buff)

	Service_Id_Buff := comm.Autocompletion([]byte(this.Service_Id), 10) // string //10 (22)
	copy(temp_buff[Body_len:], Service_Id_Buff[0:])
	Body_len += len(Service_Id_Buff)

	Fee_UserType_Buff := comm.Uint8_byte(this.Fee_UserType) // uint8 //1
	copy(temp_buff[Body_len:], Fee_UserType_Buff[0:])
	Body_len += len(Fee_UserType_Buff)

	Fee_terminal_Id_Buff := comm.Autocompletion([]byte(this.Fee_terminal_Id), 32) // string //32 (55)
	copy(temp_buff[Body_len:], Fee_terminal_Id_Buff[0:])
	Body_len += len(Fee_terminal_Id_Buff)

	Fee_terminal_type_Buff := comm.Uint8_byte(this.Fee_terminal_type) // uint8 //1
	copy(temp_buff[Body_len:], Fee_terminal_type_Buff[0:])
	Body_len += len(Fee_terminal_type_Buff)

	TP_pId_Buff := comm.Uint8_byte(this.TP_pId) // uint8 //1
	copy(temp_buff[Body_len:], TP_pId_Buff[0:])
	Body_len += len(TP_pId_Buff)

	TP_udhi_Buff := comm.Uint8_byte(this.TP_udhi) // uint8 //1
	copy(temp_buff[Body_len:], TP_udhi_Buff[0:])
	Body_len += len(TP_udhi_Buff)

	Msg_Fmt_Buff := comm.Uint8_byte(this.Msg_Fmt) // uint8 //1
	copy(temp_buff[Body_len:], Msg_Fmt_Buff[0:])
	Body_len += len(Msg_Fmt_Buff)

	Msg_src_Buff := comm.Autocompletion([]byte(this.Msg_src), 6) // string //6 (65)
	copy(temp_buff[Body_len:], Msg_src_Buff[0:])
	Body_len += len(Msg_src_Buff)

	FeeType_Buff := comm.Autocompletion([]byte(this.FeeType), 2) // string //2
	copy(temp_buff[Body_len:], FeeType_Buff[0:])
	Body_len += len(FeeType_Buff)

	FeeCode_Buff := comm.Autocompletion([]byte(this.FeeCode), 6) // string //6
	copy(temp_buff[Body_len:], FeeCode_Buff[0:])
	Body_len += len(FeeCode_Buff)

	ValId_Time_Buff := comm.Autocompletion([]byte(this.ValId_Time), 17) // string //17 (90)
	copy(temp_buff[Body_len:], ValId_Time_Buff[0:])
	Body_len += len(ValId_Time_Buff)

	At_Time_Buff := comm.Autocompletion([]byte(this.At_Time), 17) // string //17
	copy(temp_buff[Body_len:], At_Time_Buff[0:])
	Body_len += len(At_Time_Buff)

	Src_Id_Buff := comm.Autocompletion([]byte(this.Src_Id), 21) // string //21
	copy(temp_buff[Body_len:], Src_Id_Buff[0:])
	Body_len += len(Src_Id_Buff)

	DestUsr_tl_Buff := comm.Uint8_byte(this.DestUsr_tl) // uint8 //1
	copy(temp_buff[Body_len:], DestUsr_tl_Buff[0:])
	Body_len += len(DestUsr_tl_Buff)

	Dest_terminal_Id_Buff := comm.Autocompletion([]byte(this.Dest_terminal_Id), 32*int(this.DestUsr_tl)) // string //32*DestUsr_tl
	copy(temp_buff[Body_len:], Dest_terminal_Id_Buff[0:])
	Body_len += len(Dest_terminal_Id_Buff)

	Dest_terminal_type_Buff := comm.Uint8_byte(this.Dest_terminal_type) // uint8 //1
	copy(temp_buff[Body_len:], Dest_terminal_type_Buff[0:])
	Body_len += len(Dest_terminal_type_Buff)

	Msg_Length_Buff := comm.Uint8_byte(this.Msg_Length) // uint8 //1
	copy(temp_buff[Body_len:], Msg_Length_Buff[0:])
	Body_len += len(Msg_Length_Buff)

	Msg_Content_Buff := comm.Autocompletion(this.Msg_Content, int(this.Msg_Length)) // string //Msg_length string
	copy(temp_buff[Body_len:], Msg_Content_Buff[0:])
	Body_len += len(Msg_Content_Buff)

	LinkID_Buff := comm.Autocompletion([]byte(this.LinkID), 20) // string //20
	copy(temp_buff[Body_len:], LinkID_Buff[0:])
	Body_len += len(LinkID_Buff)

	body_buff = make([]byte, Body_len)
	copy(body_buff[0:], temp_buff[0:Body_len])
	return body_buff, Body_len
}
func Decode_SubmitBody_Resp(body_buff []byte) (resp *SubmitBody_Resp, err error) {
	if len(body_buff) < 12 {
		return nil, errors.New("Decode_ConnectBody_Resp:Buffer len no Enough 12")
	}
	resp = new(SubmitBody_Resp)
	resp.Msg_Id, err = comm.Byte_uint64(body_buff[0:8])
	if err != nil {
		return nil, err
	}
	resp.Result, err = comm.Byte_uint32(body_buff[8:12])
	if err != nil {
		return nil, err
	}
	return resp, nil
}
