// const.go
package cmpp

/*
Command_Id定义
*/
const (
	CMPP_CONNECT                   = uint32(0x00000001) //请求连接
	CMPP_CONNECT_RESP              = uint32(0x80000001) //请求连接应答
	CMPP_TERMINATE                 = uint32(0x00000002) //终止连接
	CMPP_TERMINATE_RESP            = uint32(0x80000002) //终止连接应答
	CMPP_SUBMIT                    = uint32(0x00000004) //提交短信
	CMPP_SUBMIT_RESP               = uint32(0x80000004) //提交短信应答
	CMPP_DELIVER                   = uint32(0x00000005) //短信下发
	CMPP_DELIVER_RESP              = uint32(0x80000005) //下发短信应答
	CMPP_QUERY                     = uint32(0x00000006) //发送短信状态查询
	CMPP_QUERY_RESP                = uint32(0x80000006) //发送短信状态查询应答
	CMPP_CANCEL                    = uint32(0x00000007) //删除短信
	CMPP_CANCEL_RESP               = uint32(0x80000007) //删除短信应答
	CMPP_ACTIVE_TEST               = uint32(0x00000008) //激活测试
	CMPP_ACTIVE_TEST_RESP          = uint32(0x80000008) //激活测试应答
	CMPP_FWD                       = uint32(0x00000009) //消息前转
	CMPP_FWD_RESP                  = uint32(0x80000009) //消息前转应答
	CMPP_MT_ROUTE                  = uint32(0x00000010) //MT路由请求
	CMPP_MT_ROUTE_RESP             = uint32(0x80000010) //MT路由请求应答
	CMPP_MO_ROUTE                  = uint32(0x00000011) //MO路由请求
	CMPP_MO_ROUTE_RESP             = uint32(0x80000011) //MO路由请求应答
	CMPP_GET_MT_ROUTE              = uint32(0x00000012) //获取MT路由请求
	CMPP_GET_MT_ROUTE_RESP         = uint32(0x80000012) //获取MT路由请求应答
	CMPP_MT_ROUTE_UPDATE           = uint32(0x00000013) //MT路由更新
	CMPP_MT_ROUTE_UPDATE_RESP      = uint32(0x80000013) //MT路由更新应答
	CMPP_MO_ROUTE_UPDATE           = uint32(0x00000014) //MO路由更新
	CMPP_MO_ROUTE_UPDATE_RESP      = uint32(0x80000014) //MO路由更新应答
	CMPP_PUSH_MT_ROUTE_UPDATE      = uint32(0x00000015) //MT路由更新
	CMPP_PUSH_MT_ROUTE_UPDATE_RESP = uint32(0x80000015) //MT路由更新应答
	CMPP_PUSH_MO_ROUTE_UPDATE      = uint32(0x00000016) //MO路由更新
	CMPP_PUSH_MO_ROUTE_UPDATE_RESP = uint32(0x80000016) //MO路由更新应答
	CMPP_GET_MO_ROUTE              = uint32(0x00000017) //获取MO路由请求
	CMPP_GET_MO_ROUTE_RESP         = uint32(0x80000017) //获取MO路由请求应答
)

/*case CMPP_FWD: //消息前转
	break
case CMPP_FWD_RESP: //消息前转应答
	break
case CMPP_MT_ROUTE: //MT路由请求
	break
case CMPP_MT_ROUTE_RESP: //MT路由请求应答
	break
case CMPP_MO_ROUTE: //MO路由请求
	break
case CMPP_MO_ROUTE_RESP: //MO路由请求应答
	break
case CMPP_GET_MT_ROUTE: //获取MT路由请求
	break
case CMPP_GET_MT_ROUTE_RESP: //获取MT路由请求应答
	break
case CMPP_MT_ROUTE_UPDATE: //MT路由更新
	break
case CMPP_MT_ROUTE_UPDATE_RESP: //MT路由更新应答
	break
case CMPP_MO_ROUTE_UPDATE: //MO路由更新
	break
case CMPP_MO_ROUTE_UPDATE_RESP: //MO路由更新应答
	break
case CMPP_PUSH_MT_ROUTE_UPDATE: //MT路由更新
	break
case CMPP_PUSH_MT_ROUTE_UPDATE_RESP: //MT路由更新应答
	break
case CMPP_PUSH_MO_ROUTE_UPDATE: //MO路由更新
	break
case CMPP_PUSH_MO_ROUTE_UPDATE_RESP: //MO路由更新应答
	break
case CMPP_GET_MO_ROUTE: //获取MO路由请求
	break
case CMPP_GET_MO_ROUTE_RESP: //获取MO路由请求应答
	break
*/
