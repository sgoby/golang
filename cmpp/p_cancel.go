/*p_cancel.go
SP向ISMG发起删除短信（CMPP_CANCEL）操作
CMPP_CANCEL操作的目的是SP通过此操作可以将已经提交给ISMG的短信删除，ISMG将以CMPP_CANCEL_RESP回应删除操作的结果。
*/
package cmpp

type CancelBody struct {
	Msg_Id uint64 //8
	/*信息标识（SP想要删除的信息标识）。*/
}

type CancelBody_Resp struct {
	Success_Id uint32 //4
	/*成功标识。
	0：成功；
	1：失败。*/
}
