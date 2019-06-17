package icmp

import (
    "fmt"
    "net"
    "sync"
    "time"
)
//parallel ping
type ParallelClient struct {
   local, remote *net.IPAddr
   count   int
   sends map[uint16]ICMP
   timeout time.Duration
   conn net.Conn
   closed chan struct{}
   results chan ICMPStatistic
   wg sync.WaitGroup
   buf []byte
}

func NewParallelClient(localAddr, remoteAddr string, count int , timeout time.Duration) (*ParallelClient, error) {
    local, err := net.ResolveIPAddr("ip", localAddr)
    if err != nil {
        return  nil, err
    }
    remote, err := net.ResolveIPAddr("ip", remoteAddr)
    if err != nil {
        return  nil, err
    }
    conn, err := net.DialIP("ip:icmp", local, remote)
    if err != nil {
        return  nil, err
    }
    conn.SetDeadline(time.Now().Add(timeout))

    client := ParallelClient{
        local: local,
        remote: remote,
        count: count,
        timeout: timeout,
        conn: conn,
        closed: make(chan struct{}),
        results: make(chan  ICMPStatistic, 1),
        buf: make([]byte, 100),
        sends: make(map[uint16]ICMP),
    }

    go client.run()
    go client.receive()

    return  &client, nil
}

func (c *ParallelClient)Close()  {
    c.closed <- struct{}{}
}


func (c *ParallelClient)run()  {
    for {
        select {
        case <- c.closed:
            c.conn.Close()
            close(c.results)
            return
        case res , ok := <-c.results:
            if ok {
                fmt.Printf(" recv %s:icmp_seq=%d ttl= %d time=%.2f ms\n", res.From, res.IcmpSeq, res.Ttl, res.Time)
                c.wg.Done()
            }else {
                return
            }
        }
    }
}

func (c *ParallelClient)receive()  {
    for {
       n, err := c.conn.Read(c.buf)
       if err != nil {
           return
       }
       c.handleResponse(c.buf[20:n])
    }
}

func (c *ParallelClient)handleResponse(b []byte)  {
     icmp := EncodeToICMP(b)
     if _, ok := c.sends[icmp.SequenceNum]; !ok {
          fmt.Println("icmp not exists , may be consumer already")
         return
     }
     delete(c.sends, icmp.Identifier)
     dur := float64(time.Now().UnixNano() - icmp.Time) / 1e6
     c.results <- ICMPStatistic{From: c.remote.String(),IcmpSeq: icmp.SequenceNum , Time: float64(dur) }
}


func (c *ParallelClient)Send()  {
    for i:=0;i < c.count ; i ++ {
        c.wg.Add(1)
        go func( i int) {
            c.write(uint16(i))
        }(i)
    }
    c.wg.Wait()
}

func (c *ParallelClient)write(n uint16)  {
	icmp := NewICMPRequest(n, time.Now().UnixNano())
	if _, ok := c.sends[n];ok {
	   fmt.Println("icmp already exists")
        return
    }
	c.sends[n] = icmp
    c.conn.Write(icmp.Bytes())
}
