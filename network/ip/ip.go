package ip

import "net"

func  GetLocalIp()([]net.IP, error)  {
    var ips []net.IP
    addrs, err := net.InterfaceAddrs()
    if err != nil {
        return  nil, err
    }
    for _, addr := range  addrs {
        if ip, ok := addr.(*net.IPNet);ok && ip.IP.To4()!= nil && !ip.IP.IsLoopback(){
            ips = append( ips, ip.IP.To4())
        }
    }
    return  ips, nil
}