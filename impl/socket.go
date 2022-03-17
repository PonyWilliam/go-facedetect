package impl

import (
	"log"
	"net"
	"strconv"
)
var conn net.Conn
var err error
var id int64
func init(){
	conn,err = net.Dial("tcp", ":54088")
	if err != nil{
		log.Fatal(err.Error() + "请先检查python是否启动")
	}
	id = 0
}
func StartListen(){
	for{
		recvdata := make([]byte,2048)
		_,err := conn.Read(recvdata)
		id,_ = strconv.ParseInt(string(recvdata),10,64)
		if err != nil{
			log.Fatal(err)
		}

	}
}
func GetSocketRes() int64{
	return id
}