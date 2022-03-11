package main

import (
	"flag"
	"fmt"
	"log"
	"os"
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
		log.Fatal("Possible Errors: ", user.Username, len(os.Args), regex_txt)
		os.Exit(1)
	}

	// Create directories.
	for _, i := range dir_list {
		os.MkdirAll(*dir+"/"+*ip+"/"+i, 0755)
		os.Chown(*dir+"/"+*ip+"/"+i, 1000, 1000)
	}

	// Append IP to /etc/hosts file.
	f, err := os.OpenFile("/etc/hosts", os.O_APPEND|os.O_WRONLY, 0600)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()
	f.WriteString(*ip + " " + *dir + "." + *platform + "\n")

	// Run Scripts
	fmt.Println("Scripts are now running...")

	go scripts.Nmap(*ip, *dir)
	scripts.GoBuster(*ip, *dir, *wordlist)
	fmt.Println("[i] Nmap: Done (Check /scans folder)")
	fmt.Println("[i] GoBuster: Done (Check /scans folder)")

}
