# go-nmap

go-nmap is a golang library to run nmap scans, parse scan results.


## Installation


```sh
go get github.com/dean2021/go-nmap
```
to install the package

```go
import "github.com/dean2021/go-nmap"
```

## Example

```go
package main

import (
	"fmt"
	"strconv"

	"github.com/dean2021/go-nmap"
)

func main() {

	n := nmap.New()

	// nmap可执行文件路径,默认不需要设置
	//m.SetSystemPath("/usr/local/nmap/bin/nmap")

	// 设置nmap扫描参数
	args := []string{"-sV", "-n", "-O", "-T4", "--open"}
	n.SetArgs(args...)

	// 设置扫描端口范围
	n.SetPorts("0-65535")

	// 设置扫描目标
	n.SetHosts("127.0.0.1")

	// 隔离扫描名单
	//n.SetExclude("127.0.0.1")

	// 开始扫描
	err := n.Run()
	if err != nil {
		fmt.Println("scanner failed:", err)
		return
	}

	// 解析扫描结果
	result, err := n.Parse()
	if err != nil {
		fmt.Println("Parse scanner result:", err)
		return
	}

	var (
		osName     string
		osAccuracy = 0
	)
	for _, host := range result.Hosts {
		if host.Status.State == "up" {

			for _, osMatch := range host.Os.OsMatches {
				tempOsAccuracy, _ := strconv.Atoi(osMatch.Accuracy)
				if tempOsAccuracy >= osAccuracy {
					osName = osMatch.Name
					osAccuracy = tempOsAccuracy
				}
			}

			ipAddr := host.Addresses[0].Addr
			for _, port := range host.Ports {
				portStr := strconv.Itoa(port.PortId)
				servicesStr := port.Service.Name
				fmt.Println(ipAddr, portStr, servicesStr, osName)
			}
		}
	}

}
```

## Thanks

https://github.com/lair-framework/go-nmap/blob/master/nmap.go