// Copyright 2020 zs. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package main

import (
	"flag"
	"fmt"
	"log"
	"runtime"
	"strconv"
	"time"

	"github.com/zs5460/ipconv"
)

var (
	maxThread     int
	fullMode      bool
	specifiedPort string
)

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

	fmt.Println(banner)

	flag.StringVar(&ip, "ip", "", "IP to be scanned, supports three formats:\n192.168.0.1 \n192.168.0.1-8 \n192.168.0.0/24")

	flag.BoolVar(&fullMode, "f", false, "Scan all TCP and UDP ports in full scan mode. The default is off. By default, only common TCP ports are scanned.")

	flag.IntVar(&maxThread, "t", 10000, "Maximum number of threads")

	flag.StringVar(&specifiedPort, "p", "", "Specific port to scan(1~65535)")

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
		if err != nil || p < 1 || p > 65535 {
			log.Fatal("Invalid Port")
		}
	}

	startTime := time.Now()

	tasks := make(chan item)
	done := make(chan int)
	go genTask(tasks, ips)
	go scan(tasks, done)
	<-done

	takes := time.Since(startTime).Truncate(time.Millisecond)
	fmt.Printf("\nScanning completed, taking %s.\n\n", takes)

}
