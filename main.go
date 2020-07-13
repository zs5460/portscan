package main

import (
	"flag"
	"fmt"
	"log"
	"net"
	"runtime"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/zs5460/ipconv"
)

var (
	maxThread     int
	fullMode      bool
	specifiedPort string
)

func isOpen(addr string) bool {
	conn, err := net.DialTimeout("tcp", addr, time.Second*1)
	if err != nil {
		// if strings.Index(err.Error(), "timeout") < 0 {
		// 	log.Println(err)
		// }
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
	defer conn.Close()
	if err != nil {
		return false
	}
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

func scan(proto string, ip string, ports []string) {
	var wg sync.WaitGroup
	limiter := make(chan struct{}, maxThread)
	result := make(chan string, 64)
	go func() {
		for s := range result {
			fmt.Println(s)
		}
	}()

	for _, p := range ports {
		wg.Add(1)
		limiter <- struct{}{}
		go func(addr string) {
			defer wg.Done()
			if proto == "udp" {
				if isOpenUDP(addr) {
					result <- fmt.Sprintf("[UDP] %22s is open", addr)
				}
			} else {
				if isOpen(addr) {
					result <- fmt.Sprintf("[TCP] %22s is open", addr)
				}
			}

			<-limiter
		}(ip + ":" + p)

	}

	wg.Wait()
	close(result)

}

// FullScan scan all ports.
func FullScan(ip string) {
	tmp := [65536]string{}
	for i := 0; i < 65536; i++ {
		tmp[i] = strconv.Itoa(i + 1)
	}
	ports := tmp[:]
	scan("tcp", ip, ports)
	scan("udp", ip, ports)
}

// QuickScan scan only common ports.
func QuickScan(ip string) {
	commonPorts := "21,22,23,25,53,80,110,135,137,138,139,443,1433,1434,1521,3306,3389,5000,5432,5632,6379,8000,8080,8081,8443,9090,10051,11211,27017"
	ports := strings.Split(commonPorts, ",")
	scan("tcp", ip, ports)
}

// SpecifiedScan scan specified port.
func SpecifiedScan(ip, port string) {
	scan("tcp", ip, strings.Split(port, ","))
}

// ScanIPS scan multiple IPs.
func ScanIPS(ips []string) {

	var wg sync.WaitGroup
	lm := 10
	if !fullMode {
		lm = 1000
	}
	limiter := make(chan struct{}, lm)
	for _, ip := range ips {
		wg.Add(1)
		limiter <- struct{}{}
		go func(ipaddr string) {
			defer wg.Done()
			if specifiedPort != "" {
				SpecifiedScan(ipaddr, specifiedPort)
			} else {
				QuickScan(ipaddr)
			}
			<-limiter
		}(ip)

	}
	wg.Wait()
}

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())

	banner := `
                      __                      	
    ____  ____  _____/ /_______________ _____ 
   / __ \/ __ \/ ___/ __/ ___/ ___/ __ \/ __ \
  / /_/ / /_/ / /  / /_(__  ) /__/ /_/ / / / /
 / .___/\____/_/   \__/____/\___/\__,_/_/ /_/ 
/_/                                           
                                     © zs5460
`

	ip := ""
	specifiedPort = ""
	fullMode = false

	fmt.Println(banner)

	flag.StringVar(&ip, "ip", "", "IP to be scanned, supports three formats:\n192.168.0.1 \n192.168.0.1-8 \n192.168.0.0/24")

	flag.BoolVar(&fullMode, "f", false, "Scan all TCP and UDP ports in full scan mode. The default is off. By default, only common TCP ports are scanned.")

	flag.IntVar(&maxThread, "t", 10000, "Maximum number of threads")

	flag.StringVar(&specifiedPort, "p", "", "Specific port to scan(0~65535)")

	flag.Parse()

	if ip == "" {
		flag.Usage()
		return
	}

	checkLimit()

	ips, err := ipconv.Parse(ip)
	if err != nil {
		log.Fatal(err)
	}

	if specifiedPort != "" {
		p, err := strconv.Atoi(specifiedPort)
		if err != nil || p < 0 || p > 65535 {
			log.Fatal("Invalid Port")
		}
	}

	startTime := time.Now()
	if len(ips) > 1 {
		if fullMode {
			fmt.Println("Multi-host mode does not support full scan")
			return
		}
		ScanIPS(ips)
	} else {
		if fullMode {
			FullScan(ip)
		} else {
			QuickScan(ip)
		}
	}
	takes := time.Since(startTime).Truncate(time.Millisecond)
	fmt.Printf("Scanning completed, taking %s.\n\n", takes)
}
