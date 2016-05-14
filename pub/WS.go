//测试封装websocket，
package pub

import (
	"flag"
	"log"
	"net/url"
	"net"
	"github.com/gorilla/websocket"
)

var addr = flag.String("addr", "192.168.1.112:80", "http service address")

type WS struct {
	conn    net.Conn
}

func (this *WS) Close(){
	this.conn.Close()
}
func (this *WS)Init() error{
	var err error
	u := url.URL{Scheme: "ws", Host: *addr, Path: "/ws"}
	log.Printf("connecting to %s", u.String())
	this, _, err = websocket.DefaultDialer.Dial(u.String(), nil)
	if err != nil {
		log.Fatal("dial:", err)
		return err
	}
	return nil
}

func (this *WS)QueryCmd(msg string) error{
	
}

func (this *WS)SendServerStat(msg string) error {
	var err error
	err = this.WriteJSON(msg)
	if err != nil {
		log.Println("write:", err)
		return err
	}
	return nil
}

func (this *WS)ReadServerStat() error {
	var err error
	_, message, err := this.ReadMessage()
	if err != nil {
		return err
	}
	log.Printf("read:%s",message)
	return nil
}

func main() {
	var (
		//c *websocket.Conn
		ws WS
	)
	/*u := url.URL{Scheme: "ws", Host: *addr, Path: "/ws"}
	log.Printf("connecting to %s", u.String())
	c, _, err = websocket.DefaultDialer.Dial(u.String(), nil)
	if err != nil {
		log.Fatal("dial:", err)
	}
	defer c.Close()*/
	ws.Init();
	defer ws.Close();
	msg := "ehehe"
	log.Printf("write:",msg)
	//conn=c;
	ws.SendServerStat(msg)
	ws.ReadServerStat()
	return
}