package convert

import (
	"encoding/json"
	"fmt"
	"testing"
)


//
type Action struct {
	CmdID     int `json:CmdID`
	Status    int `json:Status` //0上行，需回复，1下行不需回复
	ErrorCode int
	Data      interface{} `json:Data`
}

//
type User struct {
	UserID     int                  `json:UserID`
	RoomID     int                  `json:RoomID`
	Name       string               `json:Name`
	Chips      int                  `json:Chips`    //筹码
	IsFold     bool                 `json:IsFold`   //弃牌
	IsLooked   bool                 `json:IsLooked` //牌是否打开，是否已看牌
	IsOpened   bool                 `json:IsOpened`
	IsReady    bool                 `json:IsReady` //准备就续
	HasPoker   bool                 `json:HasPoker`
	MyPoker    *Poker				`json:MyPoker`
}
// \"MyPoker\":{\"Value\":3,\"Pattern\":2,\"VirtualIndex\":50,\"Avv\:[1,2,3]"}
type Poker struct {
	//大小
	Value int // VirtualIndex/4   2 ~ 14其中14表示A
	//花色
	Pattern int // VirtualIndex % 4   0:方块，1:梅花,2:红桃,3:黑桃
	//虚拟值
	VirtualIndex int //8~59
	//
	Arr []int
}

func Test_ToStruct(t *testing.T){
	testConver()

}
func testConver(){
	data := []byte("{\"CmdID\":1000,\"Status\":1,\"Data\":{\"UserID\":78899,\"Name\":\"asdf\", \"MyPoker\":{\"Value\":3,\"Pattern\":2,\"VirtualIndex\":50,\"Arr\":[1,2,3]}}}")
	mAction := new(Action)
	err := json.Unmarshal(data, &mAction)
	if err != nil {
		fmt.Println("## 无效信息", err)
		return
	}
	mUser := new(User)
	//fmt.Println(reflect.ValueOf(mAction).Elem.Kind())
	err = InterfaceToStruct(mAction.Data,&mUser)
	if err != nil{
		fmt.Println(err)
		return
	}
	fmt.Println(mAction)
	fmt.Println(mUser)
	fmt.Println(mUser.MyPoker)
}
