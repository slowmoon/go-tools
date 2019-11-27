package main

import (
	"github.com/slowmoon/go-tools/network/icmp"
)

func main()  {

	client ,err := icmp.NewICMPClient("", "35.220.222.205", 3)
	//client ,err := icmp.NewParallelClient("", "baidu.com", 10, time.Second*10)
	//client ,err := icmp.NewICMPClient("", "192.168.116.2", 3)
	if err!= nil {
		panic(err)
	}
	client.Send()
}
