/*p_deliver.go
ISMG向SP送交短信（CMPP_DELIVER）操作
CMPP_DELIVER操作的目的是ISMG把从短信中心或其它ISMG转发来的短信送交SP，SP以CMPP_DELIVER_RESP消息回应。
当ISMG向SP送交的不是状态报告时，则表示用户的上行内容即：MO
*/
package cmpp

import (
	"../comm"
	"errors"
)

type DeliverBody struct {
	Msg_Id uint64 //8
	/*信息标识。
	  生成算法如下：
	  采用64位（8字节）的整数：
	  （1）时间（格式为MMDDHHMMSS，即月日时
	  分秒）：bit64~bit39，其中
	  bit64~bit61：月份的二进制表示；
	  bit60~bit56：日的二进制表示；
	  bit55~bit51：小时的二进制表示；
	  bit50~bit45：分的二进制表示；
	  bit44~bit39：秒的二进制表示；
	  （2）短信网关代码：bit38~bit17，把短信网关的代码转换为整数填写到该字段中；
	  （3）序列号：bit16~bit1，顺序增加，步长为1，循环使用。
	  各部分如不能填满，左补零，右对齐。*/

	Dest_Id string //21
	/*目的号码。
	  SP的服务代码，一般4--6位，或者是前缀为服务代码的长号码；该号码是手机用户短消息的被叫号码。*/

	Service_Id string //10
	/*业务标识，是数字、字母和符号的组合。*/

	TP_pid uint8 //1
	/*GSM协议类型。详细解释请参考GSM03.40中的9.2.3.9。*/

	TP_udhi uint8 //1
	/*GSM协议类型。详细解释请参考GSM03.40中的9.2.3.23，仅使用1位，右对齐。*/

	Msg_Fmt uint8 //1
	/*信息格式：
	  0：ASCII串；
	  3：短信写卡操作；
	  4：二进制信息；
	  8：UCS2编码；
	  15：含GB汉字。*/

	Src_terminal_Id string //32
	/*源终端MSISDN号码（状态报告时填为CMPP_SUBMIT消息的目的终端号码）。*/

	Src_terminal_type uint8 //1
	/*源终端号码类型，0：真实号码；1：伪码。*/

	Registered_Delivery uint8 //1
	/*是否为状态报告：
	  0：非状态报告；
	  1：状态报告。*/

	Msg_Length uint8 //1
	/*消息长度，取值大于或等于0。*/

	Msg_Content []byte //Msg_Length
	/*消息内容。*/

	LinkID string //20
	/*点播业务使用的LinkID，非点播类业务的MT流程不使用该字段。*/
}

type DeliverBody_Resp struct {
	Msg_Id uint64 //8
	/*信息标识（CMPP_DELIVER中的Msg_Id字段）。*/
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
	  8: 流量控制错；
	  9~ ：其他错误。*/
}

//当ISMG向SP送交状态报告时，信息内容字段（Msg_Content）格式定义如下：
type Deliver_Msg_Content struct {
	Msg_Id uint64 //8
	/*信息标识。
	  SP提交短信（CMPP_SUBMIT）操作时，与SP相连的ISMG产生的Msg_Id。*/

	Stat string //7
	/*发送短信的应答结果，含义详见表一。SP根据该字段确定CMPP_SUBMIT消息的处理状态。*/

	Submit_time string //10
	/*YYMMDDHHMM（YY为年的后两位00-99，MM：01-12，DD：01-31，HH：00-23，MM：00-59）。*/

	Done_time string //10
	/*YYMMDDHHMM。*/

	Dest_terminal_Id string //32
	/*目的终端MSISDN号码(SP发送CMPP_SUBMIT消息的目标终端)。*/

	SMSC_sequence uint32 //4
	/*取自SMSC发送状态报告的消息体中的消息标识。*/
}

func Decode_DeliverBody(body_buff []byte) (deliver *DeliverBody, err error) {

	if len(body_buff) < 97 {
		return nil, errors.New("Decode_DeliverBody:Buffer len no Enough 97")
	}
	deliver = new(DeliverBody)
	byte_start := 0
	deliver.Msg_Id, err = comm.Byte_uint64(body_buff[byte_start : byte_start+8])
	byte_start = byte_start + 8
	if err != nil {
		return nil, err
	}

	deliver.Dest_Id = string(comm.Str_byte_end(body_buff[byte_start : byte_start+21]))
	byte_start = byte_start + 21

	deliver.Service_Id = string(comm.Str_byte_end(body_buff[byte_start : byte_start+10]))
	byte_start = byte_start + 10

	deliver.TP_pid, err = comm.Byte_uint8(body_buff[byte_start : byte_start+1])
	byte_start = byte_start + 1
	if err != nil {
		return nil, err
	}

	deliver.TP_udhi, err = comm.Byte_uint8(body_buff[byte_start : byte_start+1])
	byte_start = byte_start + 1
	if err != nil {
		return nil, err
	}

	deliver.Msg_Fmt, err = comm.Byte_uint8(body_buff[byte_start : byte_start+1])
	byte_start = byte_start + 1
	if err != nil {
		return nil, err
	}

	deliver.Src_terminal_Id = string(comm.Str_byte_end(body_buff[byte_start : byte_start+32]))
	byte_start = byte_start + 32

	deliver.Src_terminal_type, err = comm.Byte_uint8(body_buff[byte_start : byte_start+1])
	byte_start = byte_start + 1
	if err != nil {
		return nil, err
	}

	deliver.Registered_Delivery, err = comm.Byte_uint8(body_buff[byte_start : byte_start+1])
	byte_start = byte_start + 1
	if err != nil {
		return nil, err
	}

	deliver.Msg_Length, err = comm.Byte_uint8(body_buff[byte_start : byte_start+1])
	byte_start = byte_start + 1
	if err != nil {
		return nil, err
	}
	deliver.Msg_Content = body_buff[byte_start : byte_start+int(deliver.Msg_Length)]
	byte_start = byte_start + int(deliver.Msg_Length)
	//
	/*
		if deliver.Msg_Fmt == uint8(8) {
			content_uf8, err := comm.Str_ucs2_utf8(body_buff[byte_start : byte_start+int(deliver.Msg_Length)])
			if err != nil {
				return nil, err
			}
		} else {
			deliver.Msg_Content = string(body_buff[byte_start : byte_start+int(deliver.Msg_Length)])
		}
	*/
	deliver.LinkID = string(comm.Str_byte_end(body_buff[byte_start : byte_start+20]))
	//byte_start = byte_start + 20
	return deliver, nil
}
func Decode_Deliver_Msg_Content(body_buff []byte) (content *Deliver_Msg_Content, err error) {

	if len(body_buff) < 71 {
		return nil, errors.New("Decode_Deliver_Msg_Content:Buffer len no Enough 71")
	}
	content = new(Deliver_Msg_Content)
	byte_start := 0
	content.Msg_Id, err = comm.Byte_uint64(body_buff[byte_start : byte_start+8])
	byte_start = byte_start + 8
	if err != nil {
		return nil, err
	}

	content.Stat = string(comm.Str_byte_end(body_buff[byte_start : byte_start+7]))
	byte_start = byte_start + 7

	content.Submit_time = string(comm.Str_byte_end(body_buff[byte_start : byte_start+10]))
	byte_start = byte_start + 10

	content.Done_time = string(comm.Str_byte_end(body_buff[byte_start : byte_start+10]))
	byte_start = byte_start + 10

	content.Dest_terminal_Id = string(comm.Str_byte_end(body_buff[byte_start : byte_start+32]))
	byte_start = byte_start + 32

	content.SMSC_sequence, err = comm.Byte_uint32(body_buff[byte_start : byte_start+4])
	if err != nil {
		return nil, err
	}
	return content, nil
}
func Creat_DeliverBody_Resp(msgid uint64, res uint32) (body *DeliverBody_Resp) {
	body = new(DeliverBody_Resp)
	body.Msg_Id = msgid
	/*信息标识（CMPP_DELIVER中的Msg_Id字段）。*/
	body.Result = res
	return body
}
func (this *DeliverBody_Resp) Encode_DeliverBody_Resp() (body_buff []byte, totall_len int) {

	Msg_Id_buff := comm.Uint64_byte(this.Msg_Id)
	Result_buff := comm.Uint32_byte(this.Result)

	totall_len = 12
	body_buff = make([]byte, totall_len)
	copy(body_buff[0:], Msg_Id_buff[0:8])
	copy(body_buff[8:], Result_buff[0:4])
	return body_buff, totall_len
}
