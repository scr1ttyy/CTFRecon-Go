package scripts

import (
	"log"
	"os"
	"os/exec"
)

func GoBuster(ip, dir, wordlist string) {
	var gobuster_result string = dir + "/" + ip + "/scans/" + dir + "_GoBusterScan.txt"
	gobuster_path, err := exec.LookPath("gobuster")
	if err != nil {
		log.Fatal(err)
	}
	gobuster_scan := &exec.Cmd{
		Path:   gobuster_path,
		Args:   []string{gobuster_path, "dir", "-u", ip, "-w", wordlist, "-o", gobuster_result, "-t", "100"},
		Stdout: nil,
		Stderr: os.Stderr,
	}
	if err := gobuster_scan.Run(); err != nil {
		log.Fatal(err)
	}
}
