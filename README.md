# portscan

[![Build Status](https://www.travis-ci.org/zs5460/portscan.svg?branch=master)](https://www.travis-ci.org/zs5460/portscan)

A simple portscanner written in Go

## Usage

```shell
portscan -h

    ____  ____  _____/ /_______________ _____ 
   / __ \/ __ \/ ___/ __/ ___/ ___/ __ \/ __ \
  / /_/ / /_/ / /  / /_(__  ) /__/ /_/ / / / /
 / .___/\____/_/   \__/____/\___/\__,_/_/ /_/ 
/_/                                           
                                     © zs5460
Usage of portscan:
  -f    Scan all ports in full scan mode. The default is off. By default, only common ports are scan
ned.
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
