package main

import (
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strconv"
	"time"

	"github.com/PonyWilliam/go-arcsoft/RfidUtils"
	"github.com/PonyWilliam/go-arcsoft/door"
	"github.com/PonyWilliam/go-arcsoft/impl"
	"github.com/gorilla/websocket"
)

//dataurl 以及token用来访问
var dataurl string
var token string
var id int64

type Obj struct{
	name string
	nums string
	id int64
	score int64
}
type res struct{
	Code int `json:"code"`
	Msg string `json:"msg"`
	Token string `json:"token"`
}
type res2 struct{
	Code int `json:"code"`
	Data struct{Workers []Msg `json:"workers"`} `json:"data"`
}
type res3 struct{
	ID int64 `json:"id"`
	Name string `json:"product_name"`
}
type allRes struct{
	Code int `json:"code"`
	Msg string `json:"msg"`
}
type Devices struct{
	device []res3
}
type Msg struct {
	ID int `json:"ID"`
	Name string `json:"Name"`
	Nums string `json:"Nums"`
	Score int `json:"Score"`
	Telephone string `json:"Telephone"`
}

var (
	upgrader = websocket.Upgrader{
		// 允许跨域
		CheckOrigin:func(r *http.Request) bool{
			return true
		},
	}
	wsConn *websocket.Conn
	err error
	conn *impl.Connection
	data []byte
	maxindex int
	preResult bool
	baseurl string
	localPath string
	objs []Obj
	Rfid [][]byte
	count int
	flag bool
)

func ListenID(){
	
}

func init(){
	//1. 连接串口(import中已处理)，初始化RFID设备
	// door.InitSerial(("COM9"))
	// RfidUtils.InitRFID("COM8")
	// //2. 初始化连接信息、数据库等
	// dataurl = "http://weapi.dadiqq.cn"
	// token = getToken() // 获取token，用于把数据表更新
	// GetWorkerData()

	//3. 轮循RFID


	go impl.StartListen() //监听socket
	//4. 先读取人脸，RFID识别不到的话直接放行

	//5.Electron需要的websocket
}

// 获取token，有了token就可以将borrow记录写到数据库内
func getToken() string {
	fmt.Println(123)
	val := url.Values{}
	val.Set("username","admin")
	val.Set("password","admin")
	resp,err := http.PostForm(fmt.Sprintf("%swork/login", dataurl),val)
	if err != nil{
		log.Fatal("error in login")
	}
	defer resp.Body.Close()
	bs,err := ioutil.ReadAll(resp.Body)
	if err != nil{
		log.Fatal(err)
	}
	data := &res{}
	_ = json.Unmarshal(bs, &data)
	fmt.Println(data)
	if data.Code != 200{
		log.Fatal(data.Msg)
	}
	return data.Token
}

func GetWorkerData(){
	//1.数据库拉取员工
	var err error
	client := &http.Client{}
	request,err := http.NewRequest("GET",fmt.Sprintf("%swork/workers", dataurl),nil)
	if err != nil{
		log.Fatal(err)
	}
	request.Header.Add("Authorization", token) //携带token访问
	temp,_ := client.Do(request)
	response,err := ioutil.ReadAll(temp.Body)
	if err != nil{
		log.Fatal(err)
	}
	fmt.Println(string(response))
	res2 := &res2{}
	err = json.Unmarshal(response, &res2)
	if err != nil{
		log.Fatal(err)
	}
	fmt.Println(res2)
	defer temp.Body.Close()
}



//WebSocket handler
func wsHandler(w http.ResponseWriter , r *http.Request){
	//	w.Write([]byte("hello"))
	// 完成ws协议的握手操作
	if wsConn , err = upgrader.Upgrade(w,r,nil); err != nil{
		return
	}
	if conn , err = impl.InitConnection(wsConn); err != nil{
		fmt.Println(err.Error())
	}
	// 启动线程，不断发消息
	for{
		if flag{
			//获取到信息了
			wid := strconv.FormatInt(objs[maxindex].id,10)
			_ = conn.WriteMessage([]byte("w" + wid))
			_ = conn.WriteMessage([]byte("start"))//告诉客户端读取rfid信息
			for _,v := range Rfid{
				//把rfid发出去
				_ = conn.WriteMessage([]byte(hex.EncodeToString(v)))//rfid发送出去
			}
			_ = conn.WriteMessage([]byte("end"))//告诉客户端读取完毕，权利归他了！
			for {
				//开始循环读取信息
				if data , err = conn.ReadMessage();err != nil{
					log.Fatal(err)
				}
				if data != nil{
					fmt.Println(string(data))
				}
				if string(data) == "ok" {
					//o98k
					flag = false
					//通知python开始工作
					door.Send([]byte("2"))
					break //对方已处理，放行！
				}else if string(data) == "cancel"{
					//取消操作，不放行，break掉就可以
					flag = false
					//通知python开始工作
					break
				}
			}
		}else{
			//发0
			door.Send([]byte("0"))
		}
	}
}

func listen(){
	http.HandleFunc("/ws",wsHandler)
	_ = http.ListenAndServe("0.0.0.0:7777", nil)
}

func start(){
	var pre_face_id int64
	count = 0;
	for{
		if !flag{
			face_id := impl.GetSocketRes()
			if (count != 0 && pre_face_id != face_id) || face_id == -1 || face_id == 0{
				if(count > 0){
					count -= 1
				}
			}
			face_id += 2
			pre_face_id = face_id
			if count > 20{
				//可以确定是本人
				count = 0
				pre_face_id = -1
				Rfid = RfidUtils.GetNearRfid();
				r_count := 0//给定一个count计数rfid，容错
				for Rfid == nil{
					Rfid = RfidUtils.GetNearRfid();
					r_count += 1
					if r_count >= 5{
						break
					}
				}
				if r_count >=5{
					//直接开门，没有物品
					door.Send([]byte("2"))
					continue
				}
				flag = true
				//同时通知face_detect(python)暂时不需要工作了
			}
		}
	}
}

func main(){
	flag = false
	//通知python开始工作
	go listen()
	start() //主线程
	time.Sleep(time.Hour)
}