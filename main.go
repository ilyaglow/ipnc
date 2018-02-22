package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"net"
	"os"
)

func main() {
	f := flag.String("i", "", "File with CIDR notation networks, new line separated each")
	ipf := flag.String("sf", "", "File with IPs to search, new line separated each")
	ip := flag.String("s", "", "Single IP-address to search")
	flag.Parse()

	if *ipf == "" && *ip == "" {
		log.Fatal("Either -sf or -s parameter should be specified")
	}

	if *f == "" {
		log.Fatal("File with CIDRs is not specified")

	}

	nets, err := netsFromFile(*f)
	if err != nil {
		log.Fatal(err)
	}

	var ips []net.IP
	if *ipf != "" {
		ipl, err := ipsFromFile(*ipf)
		if err != nil {
			log.Fatal(err)
		}
		ips = append(ips, ipl...)
	}

	if *ip != "" {
		ips = append(ips, net.ParseIP(*ip))
	}

	for _, i := range ips {
		for _, n := range nets {
			if n.Contains(i) {
				fmt.Printf("%s has been found in %s\n", i, n.String())
			}
		}
	}
}

// netsFromFile parses file and returns a slice of net.IPNets
func netsFromFile(f string) ([]net.IPNet, error) {
	var nets []net.IPNet

	file, err := os.Open(f)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	sc := bufio.NewScanner(file)
	for sc.Scan() {
		_, n, err := net.ParseCIDR(sc.Text())
		if err != nil {
			continue
		}
		nets = append(nets, *n)
	}

	if len(nets) == 0 {
		return nil, fmt.Errorf("No valid networks found in %s", f)
	}

	return nets, nil
}

// ipsFromFile parses file and returns a slice of net.IPs
func ipsFromFile(f string) ([]net.IP, error) {
	var ips []net.IP

	file, err := os.Open(f)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	sc := bufio.NewScanner(file)
	for sc.Scan() {
		if i := net.ParseIP(sc.Text()); i != nil {
			ips = append(ips, i)
		}
	}

	if len(ips) == 0 {
		return nil, fmt.Errorf("No valid IP-addresses found in %s", f)
	}
	return ips, nil
}
