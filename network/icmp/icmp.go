package icmp

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"net"
	"time"
)

type ICMP struct {
	Type        uint8
	Code        uint8
	Checksum    uint16
	Identifier  uint16
	SequenceNum uint16
}

func NewICMPRequest() ICMP  {
	icmp :=  ICMP{
		8, 0, 0, 0, 0,
	}
	var buf bytes.Buffer
	binary.Write(&buf, binary.BigEndian, icmp)
	icmp.Checksum =  CheckSum(buf.Bytes())
	return  icmp
}

func (icmp ICMP)Bytes()[]byte  {
	var buf bytes.Buffer
	binary.Write(&buf, binary.BigEndian, icmp)
	return  buf.Bytes()
}


func CheckSum(data []byte) uint16 {
	var (
		sum    uint32
		length  = len(data)
		index  int
	)
	for length > 1 {
		sum += uint32(data[index])<<8 + uint32(data[index+1])
		index += 2
		length -= 2
	}
	if length > 0 {
		sum += uint32(data[index])
	}
	sum += sum >> 16

	return uint16(^sum)
}

type Client struct {
   localAddr ,remoteAddr  *net.IPAddr
   num  int
   resp chan ICMPStatistic
   closed chan  struct{}
   icmp ICMP
   conn *net.IPConn
}

type ICMPStatistic struct {

	From   string `json:"from"`
    IcmpSeq  int  `json:"icmp_seq"`
    Ttl int    `json:"ttl"`
    Time float64   `json:"time"`
}

func NewICMPClient(localAddr , remoteAddr string, num int) (*Client, error) {

	local , err := net.ResolveIPAddr("ip", localAddr)
	if err != nil {
		return  nil, err
	}
	remote ,err := net.ResolveIPAddr("ip", remoteAddr)
	if err != nil {
		return  nil, err
	}
	conn, err := net.DialIP("ipv4:icmp", local, remote)
	if err != nil {
		panic(err)
	}

	client := &Client{
    	localAddr: local,
    	remoteAddr: remote,
    	num: num,
    	icmp: NewICMPRequest(),
    	conn: conn,
	}
	go client.run()

  return  client, nil
}


func (c *Client)Send()  {
	for i:=0 ;i < c.num ;i ++ {
		c.resp <- c.Handle(c.conn, i)
	}
}

func (c *Client)run()  {
	for {
		select {
		case <- c.closed:
			return
		case res := <-c.resp:
			fmt.Printf("local %s recv :icmp_seq=%d ttl= %d time=%.2f\n", res.From, res.IcmpSeq, res.Ttl, res.Time)
		}
	}
}

func (c *Client)Handle(conn net.Conn, n int)  ICMPStatistic{
	res := make([]byte, 100)
	start := time.Now()
	_, err := conn.Write(c.icmp.Bytes())
	if err != nil {
		panic(err)
	}
	nl, err := conn.Read(res)
    fmt.Printf("receive response %s\n", res[:nl])
	if err != nil {
		panic(err)
	}

	dur := float64(time.Since(start).Nanoseconds())/ 1e6

	state := ICMPStatistic{
		Time: dur,
		From: c.remoteAddr.IP.String(),
		IcmpSeq: n,
	}
     return  state
}
