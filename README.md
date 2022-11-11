# portscan

![status](https://img.shields.io/github/workflow/status/zs5460/portscan/release)
[![Go Report Card](https://goreportcard.com/badge/github.com/zs5460/portscan)](https://goreportcard.com/report/github.com/zs5460/portscan)
[![codecov](https://codecov.io/gh/zs5460/portscan/branch/main/graph/badge.svg?token=b7aeunEgyb)](https://codecov.io/gh/zs5460/portscan)

A simple TCP and UDP portscanner written in Go

## Performance

* Scan all TCP and UDP ports of one host in 15 seconds.
* Scan the common ports of all hosts in a Class C subnet in 2 seconds.

## Download

[releases](https://github.com/zs5460/portscan/releases/latest)

## Usage

```shell
> portscan -h

                      __                      
    ____  ____  _____/ /_______________ _____ 
   / __ \/ __ \/ ___/ __/ ___/ ___/ __ \/ __ \
  / /_/ / /_/ / /  / /_(__  ) /__/ /_/ / / / /
 / .___/\____/_/   \__/____/\___/\__,_/_/ /_/ 
/_/                                           
                                     © zs5460
Usage of portscan:
  -f    Scan all TCP/UDP ports in full scan mode. The default is off. By default, only common TCP ports are scanned.
  -ip string
        IP to be scanned, supports three formats:
        192.168.0.1
        192.168.0.1-8
        192.168.0.0/24
  -p string
        Specific port to scan(0~65535)
  -t int
        Maximum number of threads (default 10000)
```

## Licence

Released under MIT license, see [LICENSE](LICENSE) for details.

## Special thanks

[![goland](jetbrains.png)](https://www.jetbrains.com/?from=portscan)
