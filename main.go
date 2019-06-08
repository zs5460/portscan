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
	maxThread int
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

// QuickScan scan only common ports.
func QuickScan(ip string) {
	commonPorts := "21,22,23,25,53,80,110,135,137,138,139,443,1433,1434,1521,3306,3389,5000,5432,5632,6379,8000,8080,8081,8443,9090,10051,11211,27017"
	ports := strings.Split(commonPorts, ",")
	scan(ip, ports)
}

// FullScan scan all ports.
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

	flag.StringVar(&ip, "ip", "", "IP address to scan")
	flag.BoolVar(&fullMode, "f", false, "Full scan mode scans all ports, the default is off, the default only scans common ports")
	flag.IntVar(&maxThread, "t", 10000, "Maximum number of threads")
	flag.Parse()
	if ip == "" {
		flag.Usage()
		return
	}

	startTime := time.Now()
	if fullMode {
		fmt.Println("Start a full scan...")
		FullScan(ip)
	} else {
		fmt.Println("Start a quick scan...")
		QuickScan(ip)
	}
	takes := time.Since(startTime).Truncate(time.Millisecond)
	fmt.Printf("Scan completed,it took %s.\n\n", takes)
}
