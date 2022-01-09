package main

import (
	"flag"
	"fmt"
	"os"
	"os/exec"
	"os/user"
	"regexp"

	"github.com/hambyhacks/CTFRecon-Go/scripts"
)

func main() {
	// Initialize Flags
	ip := flag.String("i", "", "IP Address.")
	dir := flag.String("d", "", "Directory to be created.")
	platform := flag.String("p", "", "Extension to be added to /etc/hosts file.")
	wordlist := flag.String("w", "common.txt", "Wordlist for Directory Busting.")
	flag.Parse()

	// Variables
	dir_list := []string{"exploit", "loot", "scans", "ss"}
	user, _ := user.Current()
	regex_txt, _ := regexp.MatchString(".txt", *wordlist)

	// GoRecon Usage
	if os.Args[1] == "-h" || os.Args[1] == "" {
		flag.PrintDefaults()
		os.Exit(1)
	}

	// Check args count, current user, and wordlist if it has `.txt` in it.
	if len(os.Args) < 1 || user.Username != "root" || !regex_txt {
		flag.PrintDefaults()
		fmt.Println("Possible Errors: ", user.Username, len(os.Args), regex_txt)
		os.Exit(1)
	}

	// Create directories.
	for _, i := range dir_list {
		os.MkdirAll(*dir+"/"+*ip+"/"+i, 0755)
	}

	// Append IP to /etc/hosts file.
	f, err := os.OpenFile("/etc/hosts", os.O_APPEND|os.O_WRONLY, 0600)
	if err != nil {
		fmt.Println(err)
	}
	defer f.Close()
	f.WriteString(*ip + " " + *dir + "." + *platform + "\n")

	// Run Scripts
	c := make(chan []byte, 2)
	Nmap(*ip, *dir, c)
	go GoBuster(*ip, *dir, *wordlist, c)

	// Change Permissions of Directories.
	perms := exec.Command("chown", "-R", "1000:1000", *dir+"/")
	if err := perms.Run(); err != nil {
		fmt.Println(err)
	}
	close(c)
}

func Nmap(ip, dir string, c chan []byte) {
	var nmap_result string = dir + "/" + ip + "/scans/" + dir + "_nmapScan.txt"
	portsToScan := scripts.PortScan(ip)
	nmap_path, err := exec.LookPath("nmap")
	if err != nil {
		fmt.Println(err)
	}
	nmap_scan := &exec.Cmd{
		Path:   nmap_path,
		Args:   []string{nmap_path, "-A", "-T4", "-p", portsToScan, "-oN", nmap_result, ip},
		Stdout: nil,
		Stderr: nil,
	}
	if err := nmap_scan.Run(); err != nil {
		fmt.Println(err)
	}

	nmap_chan, _ := nmap_scan.Output()
	c <- nmap_chan
}

func GoBuster(ip, dir, wordlist string, c chan []byte) {
	var gobuster_result string = dir + "/" + ip + "/scans/" + dir + "_GoBusterScan.txt"
	gobuster_path, err := exec.LookPath("gobuster")
	if err != nil {
		fmt.Println(err)
	}
	gobuster_scan := &exec.Cmd{
		Path:   gobuster_path,
		Args:   []string{gobuster_path, "dir", "-u", ip, "-w", wordlist, "-o", gobuster_result, "-t", "64"},
		Stdout: nil,
		Stderr: nil,
	}
	if err := gobuster_scan.Run(); err != nil {
		fmt.Println(err)
	}
	gobuster_chan, _ := gobuster_scan.Output()
	c <- gobuster_chan
}
