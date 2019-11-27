package icmp

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"math/rand"
	"net"
	"time"
)

type ICMP struct {
	Type        uint8
	Code        uint8
	Checksum    uint16
	Identifier  uint16
	SequenceNum uint16
	Time     int64
}


func NewICMPRequest(n uint16, t int64) ICMP  {
	icmp :=  ICMP{
		8, 0, 0, uint16(rand.Uint32()), n , t,
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

func EncodeToICMP(b []byte)ICMP  {

	var icmp ICMP
	 binary.Read(bytes.NewBuffer(b), binary.BigEndian,&icmp )
	return  icmp
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
    IcmpSeq  uint16  `json:"icmp_seq"`
    Ttl  byte    `json:"ttl"`
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

	conn, err := net.DialIP("ip4:icmp", local, remote)
	if err != nil {
		panic(err)
	}

   conn.SetDeadline(time.Now().Add(time.Second*20))

	client := &Client{
    	localAddr: local,
    	remoteAddr: remote,
    	num: num,
    	icmp: NewICMPRequest(2, time.Now().UnixNano()),
    	conn: conn,
    	resp: make(chan ICMPStatistic, 1),
    	closed: make(chan struct{}, 1),
	}
	go client.run()

  return  client, nil
}


func (c *Client)Send()  {
	for i:=0 ;i < c.num ;i ++ {
		c.resp <- c.Handle(uint16(i))
		time.Sleep( time.Second )
	}
}

func (c *Client)Close()  {
	c.closed <- struct{}{}
}

func (c *Client)run()  {
	for {
		select {
		case <- c.closed:
			c.conn.Close()
			close(c.resp)
			return
		case res , ok := <-c.resp:
			if ok {
				fmt.Printf(" recv %s:icmp_seq=%d ttl= %d time=%.2f ms\n", res.From, res.IcmpSeq, res.Ttl, res.Time)
			}else {
				return
			}
		}
	}
}

func (c *Client)Handle(n uint16)  ICMPStatistic{
	res := make([]byte, 1024)
	start := time.Now()
	_, err := c.conn.Write(c.icmp.Bytes())

	if err != nil {
		panic(err)
	}
	_ , err = c.conn.Read(res)
	if err != nil {
		panic(err)
	}

	dur := float64(time.Since(start).Nanoseconds())/ 1e6

	state := ICMPStatistic{
		Time: dur,
		From: c.remoteAddr.IP.String(),
		IcmpSeq: n,
		Ttl: res[8],
	}
     return  state
}
