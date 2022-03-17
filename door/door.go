package door

import (
	"io"
	"log"

	serial "github.com/tarm/goserial"
)
var s	io.ReadWriteCloser
func FatalErr(err error){
	if err != nil{
		log.Fatal(err)
	}
}
func Send(req []byte){
	_, _ = s.Write(req)
}
func InitSerial(COM string){
	var err error
	cfg := &serial.Config{Name: COM, Baud: 115200, ReadTimeout: 50}
	s,err = serial.OpenPort(cfg)
	FatalErr(err)
}