package scripts

import (
	"fmt"
	"log"
	"net"
	"os"
	"os/exec"
	"regexp"
	"strings"
	"sync"
	"time"
)

func Nmap(ip, dir string, c chan []byte) chan []byte {
	var nmap_result string = dir + "/" + ip + "/scans/" + dir + "_nmapScan.txt"
	portsToScan := PortScan(ip)

	nmap_path, err := exec.LookPath("nmap")
	if err != nil {
		log.Fatal(err)
	}

	nmap_scan := &exec.Cmd{
		Path:   nmap_path,
		Args:   []string{nmap_path, "-A", "-T4", "-Pn", "-p", portsToScan, "-oN", nmap_result, ip, "--min-rate", "5000"},
		Stdout: nil,
		Stderr: os.Stderr,
	}

	if err := nmap_scan.Run(); err != nil {
		log.Fatal(err)
	}

	nmap_out, _ := nmap_scan.Output()
	c <- nmap_out
	return c
}

func PortScan(host string) string {
	wg := sync.WaitGroup{}
	listAddr := []string{}

	for i := 1; i <= 65535; i++ {
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
	fmt.Println("[i] Open Ports: ", csvPorts)
	return csvPorts
}

func CheckTCPConnection(address string, timeout int) bool {
	_, err := net.DialTimeout("tcp", address, time.Second*time.Duration(timeout))
	return err == nil
}
