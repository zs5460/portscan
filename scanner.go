package main

import (
	"fmt"
	"net"
	"strconv"
	"sync"
	"time"
)

type item struct {
	IP    string
	Port  int
	Proto string
}

func (i item) Addr() string {
	return i.IP + ":" + strconv.Itoa(i.Port)
}

func (i item) IsOpen() bool {
	if i.Proto == "udp" {
		return isOpenUDP(i.Addr())
	}
	return isOpen(i.Addr())
}

func isOpen(addr string) bool {
	conn, err := net.DialTimeout("tcp", addr, time.Second*1)
	if err != nil {
		return false
	}
	conn.Close()
	return true
}

func isOpenUDP(addr string) bool {
	udpaddr, err := net.ResolveUDPAddr("udp4", addr)
	if err != nil {
		return false
	}
	conn, err := net.DialUDP("udp", nil, udpaddr)
	if err != nil {
		return false
	}
	defer conn.Close()
	conn.SetDeadline(time.Now().Add(time.Duration(time.Second * 1)))
	var data []byte
	switch udpaddr.Port {
	case 53:
		data = []byte("\x24\x1a\x01\x00\x00\x01\x00\x00\x00\x00\x00\x00\x03\x77\x77\x77\x06\x67\x6f\x6f\x67\x6c\x65\x03\x63\x6f\x6d\x00\x00\x01\x00\x01")
	case 123:
		data = []byte("\xe3\x00\x04\xfa\x00\x01\x00\x00\x00\x01\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\xc5\x4f\x23\x4b\x71\xb1\x52\xf3")
	case 161:
		data = []byte("\x30\x2c\x02\x01\x00\x04\x07\x70\x75\x62\x6c\x69\x63\xA0\x1E\x02\x01\x01\x02\x01\x00\x02\x01\x00\x30\x13\x30\x11\x06\x0D\x2B\x06\x01\x04\x01\x94\x78\x01\x02\x07\x03\x02\x00\x05\x00")
	default:
		data = []byte("\xff\xff\x70\x69\x65\x73\x63\x61\x6e\x6e\x65\x72\x20\x2d\x20\x40\x5f\x78\x39\x30\x5f\x5f")
	}
	_, err = conn.Write(data)
	if err != nil {
		return false
	}

	buf := make([]byte, 256)
	_, err = conn.Read(buf)
	if err != nil {
		return false
	}
	return true
}

func scan(tasks <-chan item, done chan<- int) {
	var wg sync.WaitGroup
	limiter := make(chan struct{}, maxThread)
	result := make(chan string, 1024)
	go func() {
		for s := range result {
			fmt.Println(s)
		}
	}()

	for i := range tasks {
		wg.Add(1)
		limiter <- struct{}{}
		go func(i item) {
			defer wg.Done()
			if i.IsOpen() {
				result <- fmt.Sprintf("%22s [%s] is open", i.Addr(), i.Proto)
			}
			<-limiter
		}(i)

	}

	wg.Wait()
	close(result)
	done <- 1

}

func addFullScan(tasks chan<- item, ip string) {
	for i := 1; i < 65536; i++ {
		tasks <- item{IP: ip, Port: i, Proto: "tcp"}
		tasks <- item{IP: ip, Port: i, Proto: "udp"}
	}
}

func addQuickScan(tasks chan<- item, ip string) {
	commonPorts := []int{21, 22, 23, 25, 53, 80, 110, 135, 137, 138, 139, 443, 465, 587, 1433, 1434, 1521, 3306, 3389, 3899, 4899, 5000, 5432, 5631, 5632, 6379, 8000, 8080, 8081, 8443, 9090, 10050, 10051, 11211, 27017}

	for _, p := range commonPorts {
		tasks <- item{IP: ip, Port: p, Proto: "tcp"}
	}
	tasks <- item{IP: ip, Port: 53, Proto: "udp"}
}

func addSpecifiedScan(tasks chan<- item, ip, port string) {
	if p, err := strconv.Atoi(port); err == nil {
		tasks <- item{IP: ip, Port: p, Proto: "tcp"}
		tasks <- item{IP: ip, Port: p, Proto: "udp"}
	}
}

func genTask(tasks chan item, ips []string) {
	for _, ip := range ips {
		if specifiedPort != "" {
			addSpecifiedScan(tasks, ip, specifiedPort)
		} else if fullMode {
			addFullScan(tasks, ip)
		} else {
			addQuickScan(tasks, ip)
		}
	}
	close(tasks)
}
