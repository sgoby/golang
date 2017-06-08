package main

import (
	"log"
	"net/http"
	"time"
)
//购票用户的基本结构，实际应用可以扩展
type Customer struct {
	ID              int                 //用户ID号，比如对应数据库中的账号ID
	ResponseWriter  http.ResponseWriter //句柄 或 可以定义成用户Context
	TicketLevel     int                 //想要购买票等级
	TicketAreaIndex int                 //想要购买的票区
	TicketNumber    int                 //单次够票数
	MessageChen     chan string         //消息队列
	AutoChangeArea  bool
	OrderID         int64 //购票序号，以便做数据记录
}

/*
创建一个新购票客户结构体
userID  :用户ID号
responseWriter : 句柄
number :单次够票数
areaIndex :想要购买的票区序号(A = 1,B = 2,C = 3,D = 4,E = 5)
lv :
*/
func NewCustomer(userID int, responseWriter http.ResponseWriter, number int, areaIndex int, lv int) *Customer {
	ct := new(Customer)
	ct.ID = userID //
	ct.ResponseWriter = responseWriter
	ct.TicketAreaIndex = areaIndex
	ct.TicketLevel = lv
	ct.TicketNumber = number
	ct.MessageChen = make(chan string, 2)
	if areaIndex <= 0 {
		ct.AutoChangeArea = true
	} else {
		ct.AutoChangeArea = false
	}
	return ct
}

/*
购票结果写入到消息队列
*/
func (Ct *Customer) Response(text string) {
	log.Println(text)
	Ct.MessageChen <- text
	//io.WriteString(Ct.ResponseWriter, text)
}

/*
从消息队列获取购票结果
*/
func (Ct *Customer) GetResponse() string {
	for {
		select {
		case msg := <-Ct.MessageChen:
			return msg
		case <-time.After(5 * time.Second):
			return "超时 5 秒"
		}
	}
}
