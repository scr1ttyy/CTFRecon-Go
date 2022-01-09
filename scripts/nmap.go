package scripts

import (
	"fmt"
	"net"
	"regexp"
	"strings"
	"sync"
	"time"
)

func PortScan(host string) string {
	wg := sync.WaitGroup{}
	listAddr := []string{}
	ports := 65535

	for i := 1; i <= ports; i++ {
		address := fmt.Sprintf("%s:%d", host, i)
		wg.Add(1)

		go func() {
			defer wg.Done()
			if CheckTCPConnection(address, 5) {
				listAddr = append(listAddr, address)
			}

		}()
	}
	wg.Wait()

	re := regexp.MustCompile(`(?:[0-9]+)$`)
	openPorts := []string{}
	for _, j := range listAddr {
		portNumbers := re.FindAllString(j, -1)
		openPorts = append(openPorts, portNumbers...)
	}
	csvPorts := fmt.Sprint(strings.Join(openPorts, ","))
	return csvPorts
}

func CheckTCPConnection(address string, timeout int) bool {
	_, err := net.DialTimeout("tcp", address, time.Second*time.Duration(timeout))
	return err == nil
}
