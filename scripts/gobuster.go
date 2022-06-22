package scripts

import (
	"fmt"
	"os"
	"os/exec"
)

func GoBuster(ip, dir, wordlist string) {
	var result string = dir + "/" + ip + "/scans/" + dir + "_GoBusterScan.txt"
	path, err := exec.LookPath("gobuster")
	if err != nil {
		panic(err)
	}
	gobuster_scan := &exec.Cmd{
		Path:   path,
		Args:   []string{path, "dir", "-u", ip, "-w", wordlist, "-o", result, "-t", "64"},
		Stdout: nil,
		Stderr: os.Stderr,
	}
	if err := gobuster_scan.Run(); err != nil {
		fmt.Println(err)
	}
}
