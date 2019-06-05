package main

import (
	"flag"
	"fmt"
	"net"
	"runtime"
	"strconv"
	"strings"
	"sync"
	"time"
)

var (
	wg sync.WaitGroup
)

func isOpen(addr string) bool {
	_, err := net.DialTimeout("tcp", addr, time.Second*2)
	if err != nil {
		return false
	}
	return true
}

func scan(ip string, ports []string) {
	limit := make(chan struct{}, 1000)
	for _, p := range ports {
		limit <- struct{}{}
		go func(addr string) {
			wg.Add(1)
			if isOpen(addr) {
				fmt.Printf("%s is open\n", addr)
			}
			wg.Done()
			<-limit
		}(ip + ":" + p)
	}
}

func QuickScan(ip string) {
	commonPorts := "21,22,23,25,53,80,110,135,137,138,139,443,1433,1434,1521,3306,3389,5000,5432,5632,6379,8000,8080,8081,8443,9090,10051,11211,27017"
	ports := strings.Split(commonPorts, ",")
	scan(ip, ports)
}

func FullScan(ip string) {
	tmp := [65535]string{}
	for i := 0; i < 65535; i++ {
		tmp[i] = strconv.Itoa(i + 1)
	}
	ports := tmp[:]
	scan(ip, ports)
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
	fullMode := false

	fmt.Println(banner)

	flag.StringVar(&ip, "ip", "localhost", "IP address to scan")
	flag.BoolVar(&fullMode, "f", false, "Full scan mode scans all ports, the default is off, the default only scans common ports")
	flag.Parse()

	startTime := time.Now()
	if fullMode {
		fmt.Println("Start a full scan...")
		FullScan(ip)
	} else {
		fmt.Println("Start a quick scan...")
		QuickScan(ip)
	}

	wg.Wait()

	takes := time.Since(startTime).Truncate(time.Millisecond)
	fmt.Printf("Scan completed,it took %s.\n\n", takes)
}
