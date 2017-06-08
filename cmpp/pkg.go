// pkg.go
package cmpp

import (
	"../comm"
	"errors"
)

//
type MsgPkg struct {
	Message_Header *MsgHeader //消息头(所有消息公共包头)
	Message_Body   []byte     //消息体
}

//
type MsgHeader struct {
	Total_Length uint32 //消息总长度(含消息头及消息体)
	Command_Id   uint32 //命令或响应类型
	Sequence_Id  uint32 //消息流水号,顺序累加,步长为1,循环使用（一对请求和应答消息的流水号必须相同）
}

//消息流水号
var SerialNumber uint32

func Creat_Pkg(command uint32, body []byte, sequence uint32) (pkg *MsgPkg, err error) {

	if body == nil {
		return nil, nil
	}
	pkg = new(MsgPkg)
	header := new(MsgHeader)

	header.Total_Length = uint32(len(body)) + 12
	header.Command_Id = command
	if sequence <= 0 { //主动发包
		//流水号累加
		SerialNumber += 1
		header.Sequence_Id = SerialNumber
	} else { //回复包
		header.Sequence_Id = sequence
	}

	pkg.Message_Header = header
	pkg.Message_Body = body
	//net send
	return pkg, nil
}
func (this *MsgPkg) Encode_Pkg() (pkg_buff []byte, totall_len int) {
	pkg_len := len(this.Message_Body) + 12
	pkg_buff = make([]byte, pkg_len)

	Header_buff, _ := this.Message_Header.Encode_MsgHeader()
	Body_buff := this.Message_Body

	copy(pkg_buff[0:], Header_buff[0:12])
	copy(pkg_buff[12:], Body_buff[0:])

	return pkg_buff, pkg_len
}
func Decode_Pkg(pkg_buff []byte) (pkg *MsgPkg, err error) {
	if len(pkg_buff) < 12 {
		return nil, errors.New("Decode_Pkg:Buffer len no Enough 12")
	}
	pkg = new(MsgPkg)
	pkg.Message_Header, err = Decode_MsgHeader(pkg_buff[0:12])
	if err != nil {
		return nil, err
	}
	body_len := pkg.Message_Header.Total_Length - 12
	if body_len <= 0 {
		return pkg, nil
	}
	body_buff := make([]byte, body_len)
	copy(body_buff[0:], pkg_buff[12:])
	pkg.Message_Body = body_buff

	return pkg, nil
}
func (this *MsgHeader) Encode_MsgHeader() (body_buff []byte, body_len int) {

	body_len = 12
	body_buff = make([]byte, body_len)
	Total_Length_byte := comm.Uint32_byte(this.Total_Length)
	Command_Id_byte := comm.Uint32_byte(this.Command_Id)
	Sequence_Id_byte := comm.Uint32_byte(this.Sequence_Id)

	copy(body_buff[0:], Total_Length_byte[0:4])
	copy(body_buff[4:], Command_Id_byte[0:4])
	copy(body_buff[8:], Sequence_Id_byte[0:4])

	return body_buff, body_len
}
func Decode_MsgHeader(header_buff []byte) (header *MsgHeader, err error) {
	if len(header_buff) < 12 {
		return nil, errors.New("Decode_MsgHeader:Buffer len no Enough 12")
	}
	header = new(MsgHeader)
	header.Total_Length, err = comm.Byte_uint32(header_buff[0:4])
	if err != nil {
		return nil, err
	}
	header.Command_Id, err = comm.Byte_uint32(header_buff[4:8])
	if err != nil {
		return nil, err
	}
	header.Sequence_Id, err = comm.Byte_uint32(header_buff[8:12])
	if err != nil {
		return nil, err
	}
	return header, nil
}
