package ip

import (
    "math"
    "net"
)

func  GetLocalIPv4()([]net.IP, error)  {
   return  GetLocalIpWithFilter(func(ips net.IP) bool {
        if !ips.IsLoopback() && ips.To4()!=nil {
            return  true
        }
        return false
   })
}

func GetLocalIpWithFilter(f func(net.IP)bool)([]net.IP, error) {
    var ips []net.IP
    addrs, err := net.InterfaceAddrs()
    if err != nil {
        return  nil, err
    }
    for _, addr := range  addrs {
        if ip, ok := addr.(*net.IPNet);ok &&  f(ip.IP){
            ips = append( ips, ip.IP.To4())
        }
    }
    return  ips, nil
}


func Ip4ToBytes(ip net.IP) []byte {
     if ip.To4() == nil {
         return  []byte{}
     }
     return ip.To4()[:]
}


func Bytes2Ipv4(b []byte)net.IP  {
   _ = b[3]
   return  net.IPv4(b[0], b[1], b[2], b[3])
}

//base ip and mask to get the ip range
type IP uint32
func Table(ipnet *net.IPNet)[]IP {
    var min, max  IP
    ip := ipnet.IP.To4()
    mask := ipnet.Mask

    for i:=0;i<4;i++ {
        b := IP(ip[i] & mask[i])
        min += b << (uint(3-i)*8)
    }
    ones, _ := ipnet.Mask.Size()
    max := min | IP(math.Pow(2, float64(32- ones))-1)
    var res []IP
    for i:=min ; i <=max;i ++ {
         if i & 0x000000ff ==0 {
             continue
         }
         res = append(res, i)
    }
    return res

}