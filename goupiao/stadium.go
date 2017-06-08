package main

import (
	"container/list"
	"errors"
	"fmt"
	"log"
	"sync"
	"time"
)

//队列最大等待数
const MAX_LIST_WAIT = 100000

//体育馆基本结构，实际应用可以扩展
type Stadium struct {
	AreaMap           map[int]*Area  //每个区的分布MAP
	CustomerChannel   chan *Customer //排队等待队列
	StadiumNearbyList *list.List     //排座顺序，队列
	BuyTicketRWMutex  *sync.Mutex    //锁
	RemainTickets     int            //余票
	AllAreas          int            //总区数
	AllLevelNumber    int            //总共等级数
}

//此结构体用于记录各区最优排座
type NearbyListEntity struct {
	ContinueTimes int     //连续相邻座位数
	AreaIndex     int     //区索引
	NearbySlice   []*Seat //座位数组
}

/*
创建体育馆结构体
*/
func NewStadium() (st *Stadium, err error) {
	St := new(Stadium)
	St.CustomerChannel = make(chan *Customer, MAX_LIST_WAIT)
	St.AreaMap = make(map[int]*Area)
	St.StadiumNearbyList = list.New()
	St.BuyTicketRWMutex = new(sync.Mutex)
	//创建5个32 行 82 列 的区结构体
	St.AreaMap[1], _ = NewArea(1, "A区", 32, 82)
	St.AreaMap[2], _ = NewArea(2, "B区", 32, 82)
	St.AreaMap[3], _ = NewArea(3, "C区", 32, 82)
	St.AreaMap[4], _ = NewArea(4, "D区", 32, 82)
	St.AreaMap[5], _ = NewArea(5, "E区", 32, 82)
	//初始分已卖票，基数2000左右,因为可能会随机无效虚拟坐位
	St.AreaMap[1].InitSold(600, 20, 2)
	St.AreaMap[2].InitSold(580, 20, 2)
	St.AreaMap[3].InitSold(440, 20, 2)
	St.AreaMap[4].InitSold(450, 20, 2)
	St.AreaMap[5].InitSold(550, 20, 2)
	St.AllAreas = 5
	//初始分余票信息
	St.InquireRemainTickets()
	//此处可以开启多个GO程
	go St.BuyTicketProcess(1)
	go St.BuyTicketProcess(2)
	//go St.BuyTicketProcess(3)
	//
	return St, nil
}

/*
查询各区余票总和
*/
func (St *Stadium) InquireRemainTickets() int {
	sum := 0
	str := "========= 系统余票查询 ===========\r\n"
	for i := 1; i <= 5; i++ {
		pArea, ok := St.AreaMap[i]
		if ok {
			num := pArea.GetRemainTickets()
			sum += num
			str += fmt.Sprintf("%d 区 剩余: %d 张\r\n", i, num)
		} else {
			log.Println("没找到相应区")
			//tk.Status = true
		}
	}
	str += fmt.Sprintf("总共剩余: %d 张\r\n", sum)
	log.Println(str)
	St.RemainTickets = sum
	return sum
}

/*
获取体体育馆总余票
*/
func (St *Stadium) GetRemainTickets() int {
	return St.RemainTickets
}

/*
打印详细余票信息
*/
func (St *Stadium) PrintRemainTickets() {
	str := "========= 余票列表 ===========\r\n"
	for i := 1; i <= 5; i++ {
		pArea, ok := St.AreaMap[i]
		if ok {
			pArea.PrintRemainTickets()
		} else {
			log.Println("没找到相应区")
			//tk.Status = true
		}
	}
	log.Println(str)
}

/*
购票用户进入队列，排队
*/
func (St *Stadium) QueueBuyTicket(ct *Customer) {
	select {
	case St.CustomerChannel <- ct:
		return
	case <-time.After(5 * time.Second): //超时5s
		//超过10万人时
		ct.Response("当前排人数过多，请稍候再试....")
		return
	}
}

/*
购票GO程（进程）入口
areaIndex ： Go程序号
*/
func (St *Stadium) BuyTicketProcess(areaIndex int) {
	for {
		select {
		case ct := <-St.CustomerChannel:
			St.BuyTicket(ct)
		case <-time.After(5 * time.Second): //超时5s
			log.Println("##当前空闲Go程号：", areaIndex)
		}
	}
}

/*
购票
*/
func (St *Stadium) BuyTicket(ct *Customer) {
	St.BuyTicketRWMutex.Lock()
	defer func() {
		St.BuyTicketRWMutex.Unlock()
	}()
	St.StadiumNearbyList.Init()
	//
	if ct.TicketAreaIndex > 0 && ct.TicketAreaIndex > len(St.AreaMap) {
		ct.Response("请选择有效区域或选择自动分配")
		return
	}
	if ct.TicketAreaIndex <= 0 {
		ct.TicketAreaIndex += 1
		if !St.FindSeatByArea(ct) {
			//ct.Response("请修改购票！")
			return
		}
	} else {
		if !St.FindSeatByArea(ct) {
			ct.Response("当前排队人数较多，请选稍候再试或选择其它区！")
		}
	}
}

/*
进到每个区里查找开没有卖出的座位
*/
func (St *Stadium) FindSeatByArea(ct *Customer) (status bool) {
	if ct.TicketAreaIndex > 5 {
		msgStr := fmt.Sprintf("余票不足 %d 张，请您修购票数量!", ct.TicketNumber)
		if St.StadiumNearbyList.Len() <= 0 {
			ct.Response(msgStr)
			St.PrintRemainTickets()
			return false
		}
		//
		ticks, err := St.GetBestNearbyList(ct.TicketNumber)
		if err != nil {
			ct.Response(msgStr)
			St.PrintRemainTickets()
			return false
		}
		St.OutTicket(ct, ticks)
		return true
	}
	//
	pArea, ok := St.AreaMap[ct.TicketAreaIndex]
	if !ok || pArea.GetRemainTickets() < 1 {
		if ct.AutoChangeArea {
			ct.TicketAreaIndex += 1
			return St.FindSeatByArea(ct)
		} else {
			ct.Response("请选择有效区域或选择自动分配")
			return false
		}
	}
	ticks, err := pArea.BuyTicketForNearby(ct.ID, ct.TicketNumber, ct.TicketLevel, ct.AutoChangeArea)
	if err != nil {
		log.Println(err)
		if ct.AutoChangeArea {
			if ticks != nil && len(ticks) > 0 {
				pContinueTimes := St.GetSeatContinueTime(ticks)
				St.PushNearbyList(ticks, pContinueTimes, ct.TicketAreaIndex)
			}
			ct.TicketAreaIndex += 1
			return St.FindSeatByArea(ct)
			//
		} else {
			ct.Response(fmt.Sprintf("%s", err))
			return false
		}
	} else {
		pContinueTimes := St.GetSeatContinueTime(ticks)
		if pContinueTimes >= ct.TicketNumber {
			//购买到最佳票，返回
			St.OutTicket(ct, ticks)
			return true
		} else {
			if !ct.AutoChangeArea {
				//购买到最佳票，返回
				St.OutTicket(ct, ticks)
				return true
			}
			St.PushNearbyList(ticks, pContinueTimes, ct.TicketAreaIndex)
			ct.TicketAreaIndex += 1
			return St.FindSeatByArea(ct)
		}
	}
	return false
}

/*
出票
ct : 客户
ticks ：票
*/
func (St *Stadium) OutTicket(ct *Customer, ticks []*Seat) {
	//log.Println(ticks)
	str := fmt.Sprintf("您已成功买到 %d 张票:<br>", ct.TicketNumber)
	for _, tk := range ticks {
		pArea, ok := St.AreaMap[tk.AreaIndex]
		if ok {
			pArea.SetTicketSold(tk, true)
		} else {
			log.Println("没找到相应区")
			//tk.Status = true
		}
		str += fmt.Sprintf("位置: %s, %d 等票;<br>", tk.Name, tk.Level)
	}
	remain := St.InquireRemainTickets()
	str += fmt.Sprintf("总共还有 %d 张余票:<br>", remain)
	ct.Response(str)
}

/*
从列表里取出区找到票
*/
func (St *Stadium) GetBestNearbyList(number int) (seatSlice []*Seat, err error) {
	newSeatSlice := make([]*Seat, 0)
	for {
		pElement := St.StadiumNearbyList.Front()
		if pElement == nil {
			//返回错误
			return nil, errors.New("Error Stadium 7: 余票不足!")
		}
		oNearbyListEntity := pElement.Value.(*NearbyListEntity)
		St.StadiumNearbyList.Remove(pElement)
		//
		for _, pSeat := range oNearbyListEntity.NearbySlice {
			newSeatSlice = append(newSeatSlice, pSeat)
			if len(newSeatSlice) >= number {
				return newSeatSlice, nil
			}
		}
	}
	return nil, nil
}

/*
把各区找到票放到列表里
连续次数最多的放到最前面
*/
func (St *Stadium) PushNearbyList(seatSlice []*Seat, continueTimes int, pAreaIndex int) {
	if len(seatSlice) < 1 {
		return
	}
	pNearbyListEntity := new(NearbyListEntity)
	pNearbyListEntity.ContinueTimes = continueTimes
	pNearbyListEntity.AreaIndex = pAreaIndex
	pNearbyListEntity.NearbySlice = seatSlice
	//
	if St.StadiumNearbyList.Len() < 1 {
		//往前放
		St.StadiumNearbyList.PushFront(pNearbyListEntity)
	} else {
		oNearbyListEntity := St.StadiumNearbyList.Front().Value.(*NearbyListEntity)
		if oNearbyListEntity.ContinueTimes > pNearbyListEntity.ContinueTimes {
			//往后放
			St.StadiumNearbyList.PushBack(pNearbyListEntity)
		} else {
			St.StadiumNearbyList.PushFront(pNearbyListEntity)
		}
	}
}

/*
获取坐位连续次数
*/
func (St *Stadium) GetSeatContinueTime(pSlice []*Seat) int {
	count := 1
	for k, s := range pSlice {
		if (k+1 < len(pSlice)) && ((s.ColIndex + 1) == pSlice[k+1].ColIndex) {
			count += 1
		} else {
			return count
		}
	}
	return count
}
