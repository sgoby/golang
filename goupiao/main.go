package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"strconv"
)

var mStadium *Stadium
var MaxBuyNumber = 10000
var ListenPort = ":7800"
var	userID = 0
var (
	//监听端口
	ListenPortFlag = flag.String("p", "7800", "ListenPort")
	//限制单次购票数
	MaxBuyNumberFlag = flag.String("m", "10000", "MaxBuyNumber")
)

/*
-p : 监听端口
-m : 限制单次购票数
如：goupiao.exe -p 7800 -m 100
*/
func main() {
	log.Println("=========== Star ================")
	number, err := strconv.Atoi(*MaxBuyNumberFlag)
	if err != nil && number > 0 {
		MaxBuyNumber = number
	}
	ListenPort = ":" + *ListenPortFlag
	//
	mStadium, _ = NewStadium()
	http.HandleFunc("/", Hi)
	http.HandleFunc("/buyticket", BuyTicket)
	http.ListenAndServe(ListenPort, nil)
	log.Println("=========== end ================")
}

func Hi(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Content-Type", "text/html")
	io.WriteString(w, "<a href='/buyticket?act=4'>开始购票</a>")
}

/*
购票入口
http://127.0.0.1:7800/buyticket?num=5&lv=2&area=2
num:购票数量
lv:票等级
area:区域
*/
func BuyTicket(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Content-Type", "text/html")
	query := req.URL.Query()
	//
	number := 1
	lv := 0
	areaIndex := 0
	
	//
	for k, v := range query {
		if k == "num" {
			number, _ = strconv.Atoi(v[0])
		}
		if k == "lv" {
			lv, _ = strconv.Atoi(v[0])
		}
		if k == "area" {
			areaIndex, _ = strconv.Atoi(v[0])
		}
	}
	//
	remainNumber := mStadium.GetRemainTickets()
	if remainNumber <= 0 {
		io.WriteString(w, "所有的票已售完，谢谢您的参于")
		return
	}
	//
	if number <= 0 {
		io.WriteString(w, "无效购票数量")
		return
	} else if number > MaxBuyNumber {
		msgStr := fmt.Sprintf("一次最多只能购买 %d 张票，请您修购票数量!", MaxBuyNumber)
		io.WriteString(w, msgStr)
		return
	}
	//
	if areaIndex > mStadium.AllAreas {
		io.WriteString(w, "请选择有效区域")
		return
	}
	//
	if number > remainNumber {
		msgStr := fmt.Sprintf("余票不足 %d 张，请您修购票数量!", number)
		io.WriteString(w, msgStr)
		return
	}
	//
	userID += 1 //实际可以是对应的用户账号ID
	ct := NewCustomer(userID,w, number, areaIndex, lv)
	mStadium.QueueBuyTicket(ct)
	//
	msg := ct.GetResponse()
	io.WriteString(w, msg)
}
