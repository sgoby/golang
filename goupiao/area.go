package main

import (
	"container/list"
	"errors"
	"fmt"
	"log"
	"math/rand"
	"sync"
	"time"
)

//座位结构体
type Seat struct {
	AreaIndex int    //所属区索引
	RowIndex  int    //行（排）索引
	ColIndex  int    //列(一排中座位)索引
	Name      string //名称  如：A区-16 排 19 座
	Level     int    //VIP等级
	Status    bool   //true 表示已卖出
}

//区域结构体
type Area struct {
	Index            int         //区索引 (A = 1,B = 2,C = 3,D = 4,E = 5)
	Name             string      //名称 如：A区
	SeatArray        [][]*Seat   //整个区当中所有坐位，二维数组
	AreaTotallSeat   int         //区座位索引总和（不是座位）
	RowsTotall       int         //行数(排数)
	ColsTotall       int         //列数
	AllLevelNumber   int         //坐位等级数
	BuyTicketRWMutex *sync.Mutex //锁
	NearbyList       *list.List  //记录相邻座位的列表
	RemainTickets    int         //余票
}

/*
创建一个区结构体
32 行 82 列
*/
func NewArea(index int, name string, rows int, cols int) (ar *Area, err error) {
	Ar := new(Area)
	Ar.Index = index
	Ar.Name = name
	Ar.AreaTotallSeat = rows * cols
	Ar.RowsTotall = rows
	Ar.ColsTotall = cols
	Ar.BuyTicketRWMutex = new(sync.Mutex)
	Ar.NearbyList = list.New()

	return Ar, nil
}

/*
初化已经买出的坐位
*/
func (Ar *Area) InitSold(number int, fistRowSeats int, step int) {
	if number > Ar.AreaTotallSeat {
		return
	}
	VipLevel := 0
	starpos := 0
	endpos := Ar.ColsTotall
	//
	pos := 0
	for i := 0; i < Ar.RowsTotall; i++ {
		var tmpArr []*Seat
		starpos = (Ar.ColsTotall - fistRowSeats) / 2
		endpos = starpos + fistRowSeats
		//
		fistRowSeats += step
		//
		pos = 0
		for j := 0; j < Ar.ColsTotall; j++ {
			if j < starpos || j >= endpos {
				tmpArr = append(tmpArr, nil)
				continue
			}
			pos += 1
			s := new(Seat)
			s.Name = Ar.Name + "-" + fmt.Sprintf("%d", i+1) + " 排 " + fmt.Sprintf("%d", pos) + " 座"
			s.Status = false
			s.RowIndex = i
			s.AreaIndex = Ar.Index
			//模拟8排分一个级别，实际应该以种子数据为准
			VipLevel = i / 8
			s.Level = VipLevel + 1
			//
			s.ColIndex = j
			tmpArr = append(tmpArr, s)
			Ar.RemainTickets += 1
		}
		Ar.SeatArray = append(Ar.SeatArray, tmpArr)
	}
	Ar.AllLevelNumber = VipLevel + 1
	//
	if number > 0 {
		tempList := GenerateRandomNumber(1, Ar.AreaTotallSeat, number)
		for _, index := range tempList {
			row := index / Ar.ColsTotall
			col := index % Ar.ColsTotall
			if Ar.SeatArray[row][col] == nil {
				continue
			}
			Ar.SeatArray[row][col].Status = true
			Ar.RemainTickets -= 1
		}
	}
}

/*
设置座位卖出或退回状态
 */
func (Ar *Area) SetTicketSold(pSeat *Seat, pStatus bool) {
	Ar.BuyTicketRWMutex.Lock()
	defer func() {
		Ar.BuyTicketRWMutex.Unlock()
	}()
	pSeat.Status = pStatus
	if pStatus {
		//卖出
		Ar.RemainTickets -= 1
	} else {
		//退加
		//Ar.RemainTickets += 1
	}
}

/*
打印所有座位列表（辅助调试用的）
 */
func (Ar *Area) PrintRemainTickets() {
	for _, seaRows := range Ar.SeatArray {
		for _, sea := range seaRows {
			if sea != nil && !sea.Status {
				log.Println(sea)
			}
		}
	}
}

/*
获取当前区的余票数
 */
func (Ar *Area) GetRemainTickets() int {
	Ar.BuyTicketRWMutex.Lock()
	defer func() {
		Ar.BuyTicketRWMutex.Unlock()
	}()
	return Ar.RemainTickets
}

/*
买相邻的票
*/
func (Ar *Area) BuyTicketForNearby(buyID int, number int, ticketlevel int, pAutoChangeArea bool) (arr []*Seat, err error) {
	Ar.BuyTicketRWMutex.Lock()
	log.Println(Ar.Index, "==========STAR BuyTicketForNearby=================", buyID)
	defer func() {
		log.Println(Ar.Index, "**********END BuyTicketForNearby=================", buyID)
		Ar.BuyTicketRWMutex.Unlock()
	}()
	if number < 1 {
		return nil, errors.New("Error 9: Ticket is 0")
	}
	//
	return Ar.FindNearby(number, ticketlevel)
}

/*
查找相邻的坐位
*/
func (Ar *Area) FindNearby(number int, lv int) (arr []*Seat, err error) {
	lastpos := -1
	seatSlice := make([]*Seat, 0)
	//
	Ar.NearbyList.Init()
	for _, seaRows := range Ar.SeatArray {
		lastpos = -1
		for j, sea := range seaRows {
			if sea == nil || sea.Status || j <= lastpos {
				continue
			}
			if lv > 0 && sea.Level != lv {
				continue
			}
			//log.Println(sea)
			seatSlice = append(seatSlice, sea)
			if number == 1 {
				return seatSlice, nil
			}
			for k := j + 1; k < len(seaRows); k++ {
				lastpos = k
				if seaRows[k] != nil && !seaRows[k].Status {
					//log.Println(seaRows[k])
					seatSlice = append(seatSlice, seaRows[k])
					if len(seatSlice) >= number {
						return seatSlice, nil
					}
				} else {
					Ar.PushNearbyList(seatSlice)
					seatSlice = make([]*Seat, 0)
					break
				}
			}
		}
		if len(seatSlice) < number {
			Ar.PushNearbyList(seatSlice)
			seatSlice = make([]*Seat, 0)
		} else {
			return seatSlice, nil
		}
	}
	if len(seatSlice) < number {
		Ar.PushNearbyList(seatSlice)
		seatSlice = make([]*Seat, 0)
	} else {
		return seatSlice, nil
	}
	for {
		pElement := Ar.NearbyList.Front()
		if pElement == nil {
			errStr := fmt.Sprintf("Error 7: AreaIndex = %d There is not enough tickets! len = %d", Ar.Index, len(seatSlice))
			return seatSlice, errors.New(errStr)
		}
		slice := pElement.Value.([]*Seat)
		Ar.NearbyList.Remove(pElement)
		for _, seat := range slice {
			if lv > 0 && seat.Level != lv {
				continue
			}
			seatSlice = append(seatSlice, seat)
			if len(seatSlice) >= number {
				return seatSlice, nil
			}
		}
	}
	return seatSlice, errors.New("Error 8: There is not enough tickets!")
}

/*
把相邻的座位放入列表中
相邻座位最多的放前面,类似于冒泡
 */
func (Ar *Area) PushNearbyList(seatSlice []*Seat) {
	if len(seatSlice) < 1 {
		return
	}
	//fmt.Println("PushNearbyList",seatSlice)
	if Ar.NearbyList.Len() < 1 {
		//往前放
		Ar.NearbyList.PushFront(seatSlice)
	} else {
		slice := Ar.NearbyList.Front().Value.([]*Seat)
		if len(slice) > len(seatSlice) {
			//往后放
			Ar.NearbyList.PushBack(seatSlice)
		} else {
			Ar.NearbyList.PushFront(seatSlice)
		}
	}
}

/*
生成count个[start,end)结束的不重复的随机数(用于初始化)
*/
func GenerateRandomNumber(start int, end int, count int) []int {
	if end < start || (end-start) < count {
		return nil
	}
	nums := make([]int, 0)
	//随机数生成器，加入时间戳保证每次生成的随机数不一样
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	for len(nums) < count {
		num := r.Intn((end - start)) + start
		exist := false
		for _, v := range nums {
			if v == num {
				exist = true
				break
			}
		}
		if !exist {
			nums = append(nums, num)
		}
	}
	return nums
}
