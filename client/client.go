package main

import (
	"userBind/proto"
	"github.com/gogo/protobuf/proto"
	"github.com/gorilla/websocket"
	"github.com/micro/go-micro/util/log"
	"time"
	"net/http"
	"os"
	"os/signal"
	"fmt"
	"strconv"
)

const (
	CLIENTID = 10
	USERID   = 1
)

var userIDCreator chan int

var tokenList = map[string]string {
	"1":"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1aWQiOiIxIiwibmFtZSI6IjEifQ.bmfouGIiITPyy5FgkFaqAxrdhBSVk2Ec4UydwYPkxaQ",
	"2":"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1aWQiOiIyIiwibmFtZSI6IjIifQ.SXIa-YM0S6368WOnc9Urr7WiKBoYZ0adyUXlaR1AhKM",
	"3":"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1aWQiOiIzIiwibmFtZSI6IjMifQ.ylaAVWoyUQFIxj5MhVlLBEzMz1TdDRRReb56iQgvqcE",
	"4":"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1aWQiOiI0IiwibmFtZSI6IjQifQ.quloscTeO3_UQ4kVefkkeirT7joQ6CrAAObBHUKwQXE",
	"5":"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1aWQiOiI1IiwibmFtZSI6IjUifQ.vn4U6QzuIv-gXtWB1xVJRAhwkiVKy5hKnfCOrMDhYlY",
	"6":"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1aWQiOiI2IiwibmFtZSI6IjYifQ.B0xkhz2Xf4MwrTxPGjfnlkd4m-Rh2dNrar_zKqV_ssI",
	"7":"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1aWQiOiI3IiwibmFtZSI6IjcifQ.2gnt1mqYJdo2ykaqFUS6BGe1JJ5zWpX_fRmUAnUglXw",
	"8":"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1aWQiOiI4IiwibmFtZSI6IjgifQ.tFxK2W3-pRxI7Ep1PO-RV79o0jxG-rvOD1FO_tHcqYk",
	"9":"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1aWQiOiI5IiwibmFtZSI6IjkifQ.O-lUScqpMJ4kSX55AlBAFPbsPm28HWUe5qOAEuqo0UU",
	"10":"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1aWQiOiIxMCIsIm5hbWUiOiIxMCJ9.OOMhNWh0kcuMAQb7hFSTvWOb0bXPfuX1kyCDa-2AEMI",
}

var (
	clientRes heartbeat.Request
	wsHost          = "127.0.0.1:8080"
	wsPath          = "/heartbeat"
	msgSeqId uint64 = 0
)

type Client struct {
	Host string
	Path string
}

func main() {
	userIDCreator = make(chan int, 1)
	userIDCreator <- 1
	count := 1
	go func(){
		for{
			if count>12{
				continue
			}
			fmt.Println("count: ", count)
			time.Sleep(time.Microsecond * 100)
			go msgHandler()
			count ++

		}
	}()
	time.Sleep(time.Second*1000)
	log.Log("----------->over")
}


func NewWebsocketClient(host, path string) *Client {
	return &Client{
		Host: host,
		Path: path,
	}
}

func (this *Client) SendMessage(userId uint64) error {

	// 增加一个信号监控,检测各种退出的情况,方便通知服务器断开连接
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt)

	dialer := &websocket.Dialer{
		HandshakeTimeout:time.Second * 10,
	}
	connHead := http.Header{}
	connHead.Add("UserId", strconv.Itoa(int(userId)))
	connHead.Add("UserToken", tokenList[strconv.Itoa(int(userId))])
	conn, _, err := dialer.Dial("ws://"+this.Host+this.Path, connHead)
	if err != nil {
		log.Fatal(err)
		return err
	}
	defer conn.Close() //关闭连接

	done := make(chan struct{})
	// 另外其一个goroutine处理接收消息
	go func() {
		defer close(done)
		for {
			_, message, err := conn.ReadMessage()
			if err != nil {
				log.Log("read:", err)
				return
			}
			if err := proto.Unmarshal(message, &clientRes); err != nil {
				log.Logf("proto unmarshal: %s", err)
			}
			log.Logf("recv: %v, senderId: %v, MyId: %v", clientRes, clientRes.UserId, userId)
		}
	}()
	//进行发送输入功能
	reader:= make(chan string, 1)
	reader <- "10001"
	go func(){
		for{
			time.Sleep(time.Second * 10)
			reader <- (time.Now().String() + ": " + strconv.Itoa(int(userId)))
		}
	}()
	d := ""

	for {
		select {
		case <-done:
			return nil
		case d=<-reader:
			err1 :=conn.WriteMessage(websocket.BinaryMessage, MsgAssemblerReader(d, userId))
			if err1 != nil {
				log.Logf("write close:", err1)
			} else {
				continue
			}
		case <-interrupt:
			// 发送 CloseMessage 类型的消息来通知服务器关闭连接，不然会报错CloseAbnormalClosure 1006错误
			// 等待服务器关闭连接，如果超时自动关闭.
			err := conn.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, "9999"))
			if err != nil {
				log.Fatalf("write close:", err)
				return nil
			}
			log.Fatalf("write close!")
			return nil
		}
	}
}

func getLatestUserID()uint64{
	lUID := <-userIDCreator
	nextUID := lUID+1
	userIDCreator <- nextUID
	return uint64(lUID)
}

func msgHandler() {
	clientWrapper := NewWebsocketClient(wsHost, wsPath)
	userId := getLatestUserID()

	if err := clientWrapper.SendMessage(userId); err != nil {
		log.Logf("SendMessage: errr%v", err)
	}
}

func MsgAssemblerReader(data string, userId uint64) []byte {
	msgSeqId += 1
	retPb := &heartbeat.Request{
		ClientId: CLIENTID,
		UserId:   userId,
		MsgId:    msgSeqId,
		Data:     data,
	}
	byteData, err := proto.Marshal(retPb)
	if err != nil {
		log.Fatal("pb marshaling error: ", err)
	}
	return byteData
}