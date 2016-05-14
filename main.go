package main

import (
	"bytes"
	"encoding/binary"
	"log"
	"net"
	"os"
	"io"
	"os/signal"
	"time"
	"tpub"
)

const (
	ClientID=6
	ClientName="测试服务器"
)


func main() {
	log.Println("begin dial...")
	conn, err := net.Dial("tcp", "192.168.1.112:1200")
	if err != nil {
		log.Println("dial error:", err)
		return
	}
	defer conn.Close()
	log.Println("dial ok")
	runchan(conn)
	go func(){
		for {
			var buff []byte = make([]byte,tpub.NETPACKSIZE )
			//conn.SetReadDeadline(time.Now().Add(time.Microsecond * 10000))
			n, err := conn.Read(buff)
			switch err {
				case io.EOF:
					log.Println("Err：io.EOF，Exit！")
					break;
				default:
					log.Printf("Read net pack err,Pack size：%d,err：%s\n", n,err)
					break;
			}
			if n == 0 {
				log.Println("Pack size = 0，Exit！")
				break;
			}
			log.Printf("Pack size:%d!\n", n)
		}
	}()


}
func runchan(conn net.Conn){
	ticker := time.NewTicker(time.Second * 4)
	defer ticker.Stop()
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt)
	done := make(chan struct{})
	for{
		select {
		case <-ticker.C:
			//log.Printf("定时器")
			n,err := wSysStat(conn) //每隔4秒钟发送一次系统状态
			if err != nil || n == 0 {
				log.Printf("Send Net Pack Err!,Send bytes:%d,Err:%s",n,err)
				return
			}
		case <-interrupt:
			log.Println("interrupt")
			// To cleanly close a connection, a client should send a close
			// frame and wait for the server to close the connection.
			select {
			case <-done:
			case <-time.After(time.Second):
			}
			return
		}
	}
}
func wSysStat(conn net.Conn)(int,error) {
	var (
		n int
		stat tpub.SysStat
		d    tpub.Pack
		err  error
		//buff		bytes.Buffer
	)
	stat.GetServerStatAll("eth0")
	log.Printf("stat :\n")
	//log.Printf("\t BootTime:%d-%d,%d:%d", stat.BootTime.Month(), stat.BootTime.Day(), stat.BootTime.Hour(), stat.BootTime.Minute())

	d.MsgType = int32(tpub.SendSysState)
	d.DataType = int32(tpub.SendSysState)
	d.ClientID = int32(ClientID)
	copy(d.ClientName[:],"clientname")
	copy(d.Session[:],"session")
	d.MsgLen = int32(0)
	//d.Data = stat;
	buf := new(bytes.Buffer)
	err = binary.Write(buf,binary.LittleEndian,stat)
	if err != nil {
		log.Println("binary.Write(buf,binary.LittleEndian,stat) Err:",err)
	}
	copy(d.Data[:],buf.Bytes())
	//log.Println("stat:",stat)
	//log.Println("buf:",buf)
	log.Println("d.Data:",d.Data)
	//d.Data=buf.Bytes()
	buff := new(bytes.Buffer)
	err = binary.Write(buff, binary.LittleEndian, d)
	if err != nil {
		log.Printf("binary.Write err！%s\n", err)
		return n,err
	}
	n, err = conn.Write(buff.Bytes())
	if err != nil {
		log.Printf("send net packet-SysStat err！%s\n", err)
		return n,err
	}
	log.Printf("send bytes：%d\n", n)
	return n,err
}
