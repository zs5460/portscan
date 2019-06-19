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
	conn, err := net.DialTimeout("tcp", addr, time.Second*2)
	if err != nil {
		// if strings.Index(err.Error(), "timeout") < 0 {
		// 	log.Println(err)
		// }
		return false
	}
	conn.Close()
	return true
}

func scan(ip string, ports []string) {
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

			if isOpen(addr) {
				result <- fmt.Sprintf("%s is open", addr)
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
	scan(ip, ports)
}

// QuickScan scan only common ports.
func QuickScan(ip string) {
	commonPorts := "21,22,23,25,53,80,110,135,137,138,139,443,1433,1434,1521,3306,3389,5000,5432,5632,6379,8000,8080,8081,8443,9090,10051,11211,27017"
	ports := strings.Split(commonPorts, ",")
	scan(ip, ports)
}

// SpecifiedScan scan specified port.
func SpecifiedScan(ip, port string) {
	scan(ip, strings.Split(port, ","))
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
			} else if fullMode {
				FullScan(ipaddr)
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

	flag.BoolVar(&fullMode, "f", false, "Scan all ports in full scan mode. The default is off. By default, only common ports are scanned.")

	flag.IntVar(&maxThread, "t", 10000, "Maximum number of threads")

	flag.StringVar(&specifiedPort, "p", "", "Specific port to scan(0~65535)")

	flag.Parse()

	if ip == "" {
		flag.Usage()
		return
	}

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
