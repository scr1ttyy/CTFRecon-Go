package scripts

import (
	"fmt"
	"os/exec"
)

func GoBuster(ip, dir, wordlist string, c chan []byte) {
	var gobuster_result string = dir + "/" + ip + "/scans/" + dir + "_GoBusterScan.txt"
	gobuster_path, err := exec.LookPath("gobuster")
	if err != nil {
		fmt.Println(err)
	}
	gobuster_scan := &exec.Cmd{
		Path:   gobuster_path,
		Args:   []string{gobuster_path, "dir", "-u", ip, "-w", wordlist, "-o", gobuster_result, "-t", "100"},
		Stdout: nil,
		Stderr: nil,
	}
	if err := gobuster_scan.Run(); err != nil {
		fmt.Println(err)
	}
	gobuster_chan, _ := gobuster_scan.Output()
	c <- gobuster_chan
}
